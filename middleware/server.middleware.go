// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package middleware

import (
	"sptools-backend/config"
	"sptools-backend/http_errors"
	"sptools-backend/models"

	"github.com/gin-gonic/gin"
)

func ServerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverID := c.Param("serverId")

		server, ok := config.GetServerByID(serverID)
		if !ok {
			c.Error(http_errors.ServerNotFound)
			c.Abort()
			return
		}

		c.Set("server", *server)
		c.Set("serverId", server.ID)

		c.Next()
	}
}

func CurrentServer(c *gin.Context) models.ServerConfig {
	value, exists := c.Get("server")
	if !exists {
		return models.ServerConfig{}
	}

	server, ok := value.(models.ServerConfig)
	if !ok {
		return models.ServerConfig{}
	}

	return server
}
