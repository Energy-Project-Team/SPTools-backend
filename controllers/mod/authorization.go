// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package mod_controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"sptools-backend/config"
	"sptools-backend/http_errors"
	"sptools-backend/logger"
	"sptools-backend/middleware"
	"sptools-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
)

type authSSESession struct {
	ID            string
	MinecraftUUID string
	ServerID      string
	CreatedAt     time.Time
	Events        chan map[string]any
	Done          chan struct{}
}

type discordTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type discordUserResponse struct {
	ID string `json:"id"`
}

var sseSessions = map[string]*authSSESession{}
var sseMu sync.RWMutex

func newAuthSSESession(minecraftUUID string, serverID string) *authSSESession {
	return &authSSESession{ID: uuid.New().String(), MinecraftUUID: strings.ReplaceAll(minecraftUUID, "-", ""), ServerID: serverID, CreatedAt: time.Now(), Events: make(chan map[string]any, 8), Done: make(chan struct{})}
}

func putSSESession(session *authSSESession) {
	sseMu.Lock()
	sseSessions[session.ID] = session
	sseMu.Unlock()
}

func getSSESession(id string) (*authSSESession, bool) {
	sseMu.RLock()
	session, ok := sseSessions[id]
	sseMu.RUnlock()
	return session, ok
}

func removeSSESession(id string) {
	sseMu.Lock()
	session, ok := sseSessions[id]
	if ok {
		delete(sseSessions, id)
	}
	sseMu.Unlock()
	if ok {
		select {
		case <-session.Done:
		default:
			close(session.Done)
		}
	}
}

func sendSSEEvent(c *gin.Context, event string, payload map[string]any) error {
	if event != "" {
		if _, err := fmt.Fprintf(c.Writer, "event: %s\n", event); err != nil {
			return err
		}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}

func claimString(claims map[string]interface{}, key string) string {
	value, _ := claims[key].(string)
	return value
}

func claimInt64(claims map[string]interface{}, key string) int64 {
	switch value := claims[key].(type) {
	case int64:
		return value
	case int:
		return int64(value)
	case float64:
		if math.IsNaN(value) || math.IsInf(value, 0) || math.Trunc(value) != value {
			return 0
		}
		return int64(value)
	default:
		return 0
	}
}

func (sc *ModControllers) V1SSEAuthorize(c *gin.Context) {
	server := middleware.CurrentServer(c)
	if server.ID == "" {
		c.Error(http_errors.ServerNotFound)
		return
	}

	minecraftUUID := strings.ReplaceAll(strings.TrimSpace(c.Query("minecraft_uuid")), "-", "")
	if minecraftUUID == "" {
		c.Error(models.NewCustomError(http.StatusBadRequest, "Missing minecraft_uuid", "error.client.authorization.minecraftUuid.required"))
		return
	}

	session := newAuthSSESession(minecraftUUID, server.ID)
	putSSESession(session)
	defer removeSSESession(session.ID)

	go func() {
		select {
		case <-time.After(10 * time.Minute):
			select {
			case session.Events <- map[string]any{"event": "error", "error": "Authorization timeout"}:
			default:
			}
			removeSSESession(session.ID)
		case <-session.Done:
		}
	}()

	state, err := sc.JWT.CreateToken(map[string]interface{}{"sse_id": session.ID, "minecraft_uuid": session.MinecraftUUID, "server_id": server.ID, "iat": time.Now().Unix(), "exp": time.Now().Add(10 * time.Minute).Unix()})
	if err != nil {
		logger.Error("не удалось создать state", "err", err)
		c.Error(models.NewCustomError(http.StatusInternalServerError, "Failed to create state", "error.server.authorization.state.create"))
		return
	}

	params := url.Values{"client_id": {config.GlobalConfig.DiscordClientID}, "redirect_uri": {config.GlobalConfig.DiscordRedirectURL}, "response_type": {"code"}, "scope": {"identify"}, "state": {state}}
	authURL := fmt.Sprintf("https://discord.com/api/oauth2/authorize?%s", params.Encode())

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	c.Writer.Flush()

	if err := sendSSEEvent(c, "auth", map[string]any{"auth_url": authURL}); err != nil {
		return
	}

	pingTicker := time.NewTicker(20 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-session.Done:
			return
		case event := <-session.Events:
			eventName, _ := event["event"].(string)
			delete(event, "event")
			if err := sendSSEEvent(c, eventName, event); err != nil {
				return
			}
			if eventName == "token" || eventName == "authorized" || eventName == "error" {
				return
			}
		case <-pingTicker.C:
			if _, err := fmt.Fprint(c.Writer, ": ping\n\n"); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}

func (sc *ModControllers) V1OAuthCallback(c *gin.Context) {
	code := strings.TrimSpace(c.Query("code"))
	state := strings.TrimSpace(c.Query("state"))
	errorParam := strings.TrimSpace(c.Query("error"))

	if code == "" && errorParam == "" {
		logger.Warn("отсутствует authorization code")
		c.Error(http_errors.AuthorizationCallbackCode)
		return
	}
	if state == "" {
		logger.Warn("отсутствует state")
		c.Error(http_errors.AuthorizationCallbackState)
		return
	}

	claims, err := sc.JWT.ValidateToken(state)
	if err != nil {
		logger.Warn("недопустимый state", "err", err)
		c.Error(http_errors.AuthorizationInvalidState)
		return
	}

	sseID := claimString(claims, "sse_id")
	minecraftUUID := strings.ReplaceAll(claimString(claims, "minecraft_uuid"), "-", "")
	serverID := claimString(claims, "server_id")
	if sseID == "" || minecraftUUID == "" || serverID == "" || claimInt64(claims, "exp") < time.Now().Unix() {
		logger.Warn("state не содержит обязательные поля", "sseID", sseID, "minecraftUUID", minecraftUUID, "serverID", serverID)
		c.Error(http_errors.AuthorizationInvalidState)
		return
	}

	session, exists := getSSESession(sseID)
	if !exists {
		logger.Warn("Не найдена SSE-сессия", "sseID", sseID)
		c.Error(http_errors.AuthorizationSSENotFound)
		return
	}
	defer removeSSESession(sseID)

	if errorParam == "access_denied" {
		logger.Warn("пользователю отказано в авторизации")
		select {
		case session.Events <- map[string]any{"event": "error", "error": "Authorization denied by user"}:
		default:
		}
		c.Error(http_errors.AuthorizationAccessDenied)
		return
	}

	client := req.C().SetTimeout(15*time.Second).SetCommonHeader("User-Agent", "Energy-Project-Team/SPTools-backend/"+config.GlobalConfig.AppVersion)
	var tokenData discordTokenResponse
	resp, err := client.R().SetFormData(map[string]string{"client_id": config.GlobalConfig.DiscordClientID, "client_secret": config.GlobalConfig.DiscordClientSecret, "grant_type": "authorization_code", "code": code, "redirect_uri": config.GlobalConfig.DiscordRedirectURL, "scope": "identify"}).SetSuccessResult(&tokenData).Post("https://discord.com/api/oauth2/token")
	if err != nil {
		logger.Error("не удалось получить токен Discord", "err", err)
		c.Error(http_errors.AuthorizationTokenRequest)
		return
	}
	if !resp.IsSuccessState() || tokenData.AccessToken == "" {
		logger.Error("не удалось получить токен Discord", "status", resp.StatusCode, "body", resp.String())
		c.Error(http_errors.AuthorizationTokenRequest)
		return
	}

	var discordUser discordUserResponse
	resp, err = client.R().SetHeader("Authorization", "Bearer "+tokenData.AccessToken).SetSuccessResult(&discordUser).Get("https://discord.com/api/users/@me")
	if err != nil {
		logger.Error("не удалось получить пользовательские данные Discord", "err", err)
		c.Error(http_errors.AuthorizationUserDataRequest)
		return
	}
	if !resp.IsSuccessState() || discordUser.ID == "" {
		logger.Error("не удалось получить данные пользователя Discord", "status", resp.StatusCode, "body", resp.String())
		c.Error(http_errors.AuthorizationUserDataRequest)
		return
	}

	spworldsClient, ok := config.GetSPworldsClientByID(serverID)
	if !ok {
		c.Error(http_errors.ServerNotFound)
		return
	}
	spworldsUser, err := spworldsClient.User(discordUser.ID)
	if err != nil || strings.ReplaceAll(*spworldsUser.UUID, "-", "") != minecraftUUID {
		logger.Warn("игрок не найден на SP Worlds или UUID не совпал", "serverId", serverID, "discordID", discordUser.ID, "err", err, "spworldsUser", spworldsUser)
		c.Error(http_errors.AuthorizationSPWorldsMismatch)
		return
	}

	minecraftProfile, err := sc.PlayerDB.GetProfileByUUID(minecraftUUID)
	if err != nil || minecraftProfile.Username == "" || !strings.EqualFold(minecraftProfile.Username, *spworldsUser.Username) {
		logger.Warn("несоответствие ников", "serverId", serverID, "minecraftUUID", minecraftUUID, "expected", spworldsUser.Username, "actual", minecraftProfile.Username, "err", err)
		c.Error(http_errors.AuthorizationSPWorldsMismatch)
		return
	}

	user, err := sc.DataBase.UpsertUserLogin(c, minecraftUUID, discordUser.ID, serverID)
	if err != nil {
		logger.Error("не удалось создать или обновить пользователя", "serverId", serverID, "discordID", discordUser.ID, "minecraftUUID", minecraftUUID, "err", err)
		c.Error(http_errors.AuthorizationCreateUser)
		return
	}
	if user.Status != "active" {
		select {
		case session.Events <- map[string]any{"event": "error", "error": "User is not active"}:
		default:
		}
		c.Error(http_errors.AuthorizationUserNotFound)
		return
	}

	now := time.Now()
	ttl := time.Duration(config.GlobalConfig.JWTTokenExpireMinutes) * time.Minute
	token, err := sc.JWT.CreateToken(map[string]interface{}{"sub": user.ID, "ver": user.TokenVersion, "iat": now.Unix(), "exp": now.Add(ttl).Unix()})
	if err != nil {
		logger.Error("не удалось сгенерировать JWT", "userUUID", user.ID, "err", err)
		c.Error(http_errors.AuthorizationCreateToken)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jeff", token, int(ttl.Seconds()), "/", "", config.GlobalConfig.Node == models.Production, true)

	select {
	case session.Events <- map[string]any{"event": "authorized", "serverId": serverID, "userId": user.ID}:
	default:
		logger.Error("не удалось отправить authorized в SSE: канал переполнен")
		c.Error(models.NewCustomError(http.StatusInternalServerError, "Failed to send authorization event", "error.server.authorization.sse.sendToken"))
		return
	}

	logger.Info("Авторизация прошла успешно", "serverId", serverID, "discordID", discordUser.ID, "minecraftUUID", minecraftUUID, "userUUID", user.ID)
	c.JSON(http.StatusOK, models.NewCustomResponse(http.StatusOK, "Authorization successful"))
}

func (sc *ModControllers) V1CheckAuthorization(c *gin.Context) {
	userData, _ := c.Get("user")
	user, _ := userData.(models.User)
	c.JSON(http.StatusOK, user)
}

func (sc *ModControllers) V1LogoutAll(c *gin.Context) {
	userData, _ := c.Get("user")
	user, _ := userData.(models.User)
	if err := sc.DataBase.IncrementUserTokenVersion(c, user.ID); err != nil {
		logger.Error("не удалось выполнить logout", "userUUID", user.ID, "err", err)
		c.Error(http_errors.AuthorizationFindUser)
		return
	}
	c.SetCookie("jeff", "", -1, "/", "", config.GlobalConfig.Node == models.Production, true)
	logger.Info("logout выполнен", "userUUID", user.ID)
	c.JSON(http.StatusOK, models.NewCustomResponse(http.StatusOK, "You have successfully closed all sessions."))
}
