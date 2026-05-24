// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package general_controllers

import (
	"sptools-backend/controllers"
	"sptools-backend/database"
	"sptools-backend/services"

	"github.com/gin-gonic/gin"
)

type GeneralControllers struct {
	Engine   *gin.Engine
	DataBase *database.DataBase
	PlayerDB *services.PlayerDBService
	Calling  *services.CallingService
	JWT      *services.JWTService
}

func SetupGeneralControllers(ic *controllers.InitControllers) *GeneralControllers {
	return &GeneralControllers{
		Engine:   ic.Engine,
		DataBase: ic.DataBase,
		PlayerDB: ic.PlayerDB,
		Calling:  ic.Calling,
		JWT:      ic.JWT,
	}
}
