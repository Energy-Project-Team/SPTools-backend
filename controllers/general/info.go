// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package general_controllers

import (
	"time"

	"github.com/gin-gonic/gin"

	"sptools-backend/config"
	"sptools-backend/models"
)

func (ic *GeneralControllers) Info(c *gin.Context) {
	c.JSON(200, models.InfoResponse{
		Node:            config.GlobalConfig.Node,
		AppVersion:      config.GlobalConfig.AppVersion,
		Service:         "sptools-backend",
		DiscordClientID: config.GlobalConfig.DiscordClientID,
		StartAt:         config.GlobalConfig.StartAt.Format(time.RFC3339),
		Uptime:          float64(time.Since(config.GlobalConfig.StartAt)) / float64(time.Second),
	})
}
