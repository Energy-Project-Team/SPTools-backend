// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package general_controllers

import (
	"net/http"

	"sptools-backend/config"

	"github.com/gin-gonic/gin"
)

func (ic *GeneralControllers) Servers(c *gin.Context) {
	servers := config.GetEnabledServers()
	c.JSON(http.StatusOK, servers)
}