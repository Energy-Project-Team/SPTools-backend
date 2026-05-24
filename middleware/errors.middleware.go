// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package middleware

import (
	"fmt"
	"net/http"

	"sptools-backend/logger"
	"sptools-backend/models"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func RouteNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Message:    fmt.Sprintf("Route %s:%s not found", c.Request.Method, c.Request.URL.Path),
		Error:      http.StatusText(http.StatusNotFound),
		StatusCode: http.StatusNotFound,
		Code:       "error.client.route.notFound",
	})
}

func ErrorProcessing(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err
		if customErr, ok := err.(*models.CustomErrorResponse); ok {
			user, _ := c.Get("user")
			logger.Warn("Ошибка в обработке запроса", "errMessage", customErr.Message, "errCode", customErr.Code, "err", err, "path", c.Request.URL.Path, "user", user)
			c.JSON(customErr.StatusCode, models.ErrorResponse{
				Message:    customErr.Message,
				Error:      http.StatusText(customErr.StatusCode),
				StatusCode: customErr.StatusCode,
				Code:       customErr.Code,
			})
		} else {
			bgRed := color.New(color.BgHiRed)
			bold := color.New(color.Bold).SprintFunc()
			italic := color.New(color.Italic).SprintFunc()
			bgRed.Println(bold("ВОЗНИКЛА НЕИЗВЕСТНАЯ ОШИБКА:"), italic(err))
			user, _ := c.Get("user")
			logger.Error("ВОЗНИКЛА НЕИЗВЕСТНАЯ ОШИБКА!", "err", err, "path", c.Request.URL.Path, "user", user)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Message:    "Unknown internal server error",
				Error:      http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       "error.server.unknown",
			})
		}
		c.Abort()
		return
	}
}
