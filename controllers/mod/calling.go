// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package mod_controllers

import (
	"fmt"
	"net/http"
	"strings"

	"sptools-backend/config"
	"sptools-backend/http_errors"
	"sptools-backend/logger"
	"sptools-backend/middleware"
	"sptools-backend/models"

	"github.com/gin-gonic/gin"
)

func (sc *ModControllers) V1Calling(c *gin.Context) {
	server := middleware.CurrentServer(c)

	if server.ID == "" {
		c.Error(http_errors.ServerNotFound)
		return
	}

	if !server.Features.Calling {
		c.Error(http_errors.FunctionDisabled)
		return
	}

	userData, _ := c.Get("user")
	user, _ := userData.(models.User)

	var input struct {
		Service     string `json:"service"`
		Coordinates string `json:"coordinates"`
		Comment     string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(http_errors.InvalidBody)
		return
	}

	input.Service = strings.TrimSpace(input.Service)
	input.Comment = strings.TrimSpace(input.Comment)

	if (input.Service == "") || (input.Comment == "") {
		c.Error(http_errors.InvalidBody)
		return
	}

	callingWebhook, found := config.GetServerCallingWebhook(server.ID, input.Service)
	if !found {
		availableServices := getAvailableCallingServices(server.ID)

		logger.Warn("неправильный сервис вызова", "serverId", server.ID, "service", input.Service, "availableServices", strings.Join(availableServices, ", "))

		c.Error(models.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Available services: %s", strings.Join(availableServices, ", ")), "error.client.calling.service.invalid"))
		return
	}

	profile, err := sc.PlayerDB.GetProfileByUUID(user.MinecraftUUID)
	if err != nil {
		logger.Error("ошибка при получении псевдонима по UUID", "serverId", server.ID, "minecraftUUID", user.MinecraftUUID, "err", err)

		c.Error(http_errors.CallingGetNickname)
		return
	}

	responseData, statusCode, err := sc.Calling.Calling(callingWebhook.Webhook, profile.Username, user.DiscordID, callingWebhook.RoleID, input.Comment, input.Coordinates)
	if err != nil {
		logger.Error("ошибка при отправке Discord webhook", "serverId", server.ID, "service", input.Service, "err", err)

		c.Error(http_errors.CallingSendWebhook)
		return
	}

	if statusCode/100 == 2 {
		logger.Info("Вызов был успешно завершён", "serverId", server.ID, "service", input.Service, "discordID", user.DiscordID, "nickname", profile.Username, "uuid", user.MinecraftUUID)

		c.JSON(http.StatusOK, models.CallingSuccessResponse{
			Message:       "The call was successfully completed",
			Service:       input.Service,
			DiscordID:     user.DiscordID,
			Nickname:      profile.Username,
			MinecraftUUID: user.MinecraftUUID,
		})
		return
	}

	logger.Error("не удалось сделать вызов сотрудника", "serverId", server.ID, "statusCode", statusCode, "service", input.Service, "discordID", user.DiscordID, "nickname", profile.Username, "uuid", user.MinecraftUUID)

	c.JSON(http.StatusBadGateway, models.CallingFailureResponse{
		Message:       "Failed to make a call",
		Service:       input.Service,
		DiscordID:     user.DiscordID,
		Nickname:      profile.Username,
		MinecraftUUID: user.MinecraftUUID,
		DiscordResponse: map[string]interface{}{
			"status_code": statusCode,
			"json":        responseData,
		},
	})
}

func getAvailableCallingServices(serverID string) []string {
	webhooks := config.GetServerCallingWebhooks(serverID)

	services := make([]string, 0, len(webhooks))
	for _, webhook := range webhooks {
		services = append(services, webhook.Service)
	}

	return services
}
