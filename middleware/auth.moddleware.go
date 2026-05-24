// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package middleware

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"sptools-backend/http_errors"
	"sptools-backend/logger"
	"sptools-backend/models"
	"sptools-backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func claimToInt64(v interface{}) (int64, bool) {
	switch t := v.(type) {
	case int64:
		return t, true
	case int:
		return int64(t), true
	case float64:
		if math.IsNaN(t) || math.IsInf(t, 0) || math.Trunc(t) != t {
			return 0, false
		}
		return int64(t), true
	default:
		return 0, false
	}
}

func tokenFromRequest(c *gin.Context) (string, bool) {
	if token, err := c.Cookie("jeff"); err == nil && token != "" {
		return token, true
	}
	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		return strings.TrimSpace(authorization[7:]), true
	}
	return "", false
}

func findUserByToken(c *gin.Context, mongoDB *mongo.Database, jwt *services.JWTService) (models.User, error) {
	token, ok := tokenFromRequest(c)
	if !ok {
		return models.User{}, http_errors.AuthorizationMissingToken
	}

	claims, err := jwt.ValidateToken(token)
	if err != nil {
		logger.Error("ошибка при проверке токена", "err", err)
		return models.User{}, http_errors.AuthorizationInvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		logger.Error("ошибка аутентификации: неверный sub")
		return models.User{}, http_errors.AuthorizationInvalidToken
	}

	tokenVer, ok := claimToInt64(claims["ver"])
	if !ok {
		logger.Error("ошибка аутентификации: неверный ver")
		return models.User{}, http_errors.AuthorizationInvalidToken
	}

	var user models.User
	err = mongoDB.Collection("Users").FindOne(c, bson.M{"id": userID, "tokenVersion": tokenVer}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		logger.Error("пользователь не найден или tokenVersion не совпал", "userID", userID, "ver", tokenVer)
		return models.User{}, http_errors.AuthorizationUserNotFound
	}
	if err != nil {
		logger.Error("ошибка поиска пользователя в БД", "err", err, "userId", userID, "tokenVer", tokenVer)
		return models.User{}, http_errors.AuthorizationFindUser
	}

	if user.IsBanned {
		if user.BannedAt != nil && user.BanReason != nil {
			return models.User{}, models.NewCustomError(http.StatusForbidden, "Account has been banned on "+user.BannedAt.Format(time.RFC822)+" for the following reasons \""+*user.BanReason+"\".", "error.client.ban.atAndReason")
		}
		if user.BannedAt != nil {
			return models.User{}, models.NewCustomError(http.StatusForbidden, "Account has been banned on "+user.BannedAt.Format(time.RFC822)+".", "error.client.ban.at")
		}
		if user.BanReason != nil {
			return models.User{}, models.NewCustomError(http.StatusForbidden, "Account has been banned for the following reasons \""+*user.BanReason+"\".", "error.client.ban")
		}
		return models.User{}, models.NewCustomError(http.StatusForbidden, "Account has been banned.", "error.client.ban")
	}

	if serverID := strings.TrimSpace(c.Param("serverId")); serverID != "" && !user.HasServer(serverID) {
		return models.User{}, models.NewCustomError(http.StatusForbidden, "User is not authorized on this server", "error.client.authorization.server.not_authorized")
	}

	return user, nil
}

func AuthMiddleware(mongoDB *mongo.Database, jwt *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := findUserByToken(c, mongoDB, jwt)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			logger.Error("ошибка при преобразовании пользователя в JSON", "err", err)
			c.Abort()
			return
		}

		logger.Info(fmt.Sprintf("Приватный запрос.\nПользователь: %s.\nХэндлер: %s.\nПуть: %s\n", string(userJSON), c.HandlerName(), c.Request.URL.Path))
		c.Set("user", user)
		c.Next()
	}
}

func NoAuthMiddleware(mongoDB *mongo.Database, jwt *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := findUserByToken(c, mongoDB, jwt)
		if err != nil {
			c.Next()
			return
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			c.Next()
			return
		}

		logger.Info(fmt.Sprintf("Публичный запрос.\nПользователь: %s.\nХэндлер: %s.\nПуть: %s\n", string(userJSON), c.HandlerName(), c.Request.URL.Path))
		c.Set("user", user)
		c.Next()
	}
}
