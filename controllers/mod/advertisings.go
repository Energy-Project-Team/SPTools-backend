// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package mod_controllers

import (
	"net/http"

	"sptools-backend/http_errors"
	"sptools-backend/logger"
	"sptools-backend/middleware"
	"sptools-backend/models"

	"github.com/gin-gonic/gin"
)

func (sc *ModControllers) Advertisings(c *gin.Context) {
	server := middleware.CurrentServer(c)
	if server.ID == "" {
		c.Error(http_errors.ServerNotFound)
		return
	}
	if !server.Features.Advertisings {
		c.Error(http_errors.FunctionDisabled)
		return
	}

	items, err := sc.DataBase.Advertisings(c, server.ID)
	if err != nil {
		logger.Error("Ошибка при получении рекламных объявлений из БД", "serverId", server.ID, "err", err)
		c.Error(http_errors.AdvertisingsGetData)
		return
	}

	logger.Info("Рекламные объявления успешно получены", "serverId", server.ID, "count", len(items))
	if len(items) == 0 {
		items = make([]models.AdvertisingItem, 0)
	}
	c.JSON(http.StatusOK, items)
}
