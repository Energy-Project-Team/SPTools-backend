// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	allowedOrigins := []string{
		"http://localhost:4173",
		"https://sptools.energyproject.dev",
	}
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
