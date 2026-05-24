// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package controllers

import (
	"sptools-backend/database"
	"sptools-backend/services"

	"github.com/gin-gonic/gin"
)

type InitControllers struct {
	Engine   *gin.Engine
	DataBase *database.DataBase
	PlayerDB *services.PlayerDBService
	Calling  *services.CallingService
	JWT      *services.JWTService
}

func SetupInitControllers(engine *gin.Engine, database *database.DataBase) *InitControllers {
	return &InitControllers{
		Engine:   engine,
		DataBase: database,
		PlayerDB: services.NewPlayerDBService(),
		Calling:  services.NewCallingService(),
		JWT:      services.NewJWTService(),
	}
}
