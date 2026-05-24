// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package mod_controllers

import (
	"sptools-backend/controllers"
	"sptools-backend/database"
	"sptools-backend/services"

	"github.com/gin-gonic/gin"
)

type ModControllers struct {
	Engine   *gin.Engine
	DataBase *database.DataBase
	PlayerDB *services.PlayerDBService
	Calling  *services.CallingService
	JWT      *services.JWTService
}

func SetupModControllers(ic *controllers.InitControllers) *ModControllers {
	return &ModControllers{
		Engine:   ic.Engine,
		DataBase: ic.DataBase,
		PlayerDB: ic.PlayerDB,
		Calling:  ic.Calling,
		JWT:      ic.JWT,
	}
}
