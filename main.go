// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"context"
	"fmt"

	"sptools-backend/config"
	"sptools-backend/controllers"
	general_controllers "sptools-backend/controllers/general"
	mod_controllers "sptools-backend/controllers/mod"
	"sptools-backend/database"
	"sptools-backend/logger"
	"sptools-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jedib0t/go-pretty/v6/table"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Database
var DB *database.DataBase

func init() {
	const appVersion = "1.0.0"
	logger.InitLogger(logger.LoggerConfig{AddSource: true, NoColor: false, LogDir: "logs"})

	if err := config.LoadConfig(appVersion); err != nil {
		logger.Fatal("Не удалось загрузить конфигурацию", "err", err)
	}

	logger.Info("Создаю подключение к MongoDB")
	clientOptions := options.Client().ApplyURI("mongodb://" + config.GlobalConfig.MongoInitDBRootUsername + ":" + config.GlobalConfig.MongoInitDBRootPassword + "@" + config.GlobalConfig.MongoInitDBHost + "/")
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Fatal("Не удалось подключиться к MongoDB", "err", err)
	}
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		logger.Fatal("Не удалось пингануть MongoDB", "err", err)
	}
	mongoDB = mongoClient.Database("sptools-" + string(config.GlobalConfig.Node))
}

func main() {
	t := table.NewWriter()
	t.SetCaption("Информация о приложении")
	t.AppendHeader(table.Row{"№П/п", "Параметр", "Значение"})
	t.AppendRows([]table.Row{{1, "Версия", config.GlobalConfig.AppVersion}, {2, "Дата и время запуска", config.GlobalConfig.StartAt.String()}, {3, "Сервис", "SPTools-backend"}, {4, "Режим работы", string(config.GlobalConfig.Node)}, {5, "Имя базы данных (MongoDB)", mongoDB.Name()}})
	fmt.Println(t.Render())

	engine := gin.Default()
	DB = database.SetupDataBase(engine, mongoDB)
	if err := DB.CreateIndexes(context.Background()); err != nil {
		logger.Fatal("Не удалось создать индексы", "err", err)
	}

	initControllers := controllers.SetupInitControllers(engine, DB)
	general := general_controllers.SetupGeneralControllers(initControllers)
	mod := mod_controllers.SetupModControllers(initControllers)

	engine.NoRoute(middleware.RouteNotFound)
	engine.Use(middleware.ErrorProcessing)
	engine.Use(middleware.CORSMiddleware)

	publicOptional := engine.Group("/api/")

	private := engine.Group("/api")
	private.Use(middleware.AuthMiddleware(mongoDB, initControllers.JWT))

	publicOptional.GET("/info", general.Info)
	publicOptional.GET("/servers", general.Servers)
	publicOptional.GET("/v1/versions", general.V1Versions)
	publicOptional.GET("/v1/version/latest", general.V1LatestVersion)

	modPublic := publicOptional.Group("/mod/:serverId/v1").Use(middleware.ServerMiddleware())
	modPublic.GET("/authorization/sse", mod.V1SSEAuthorize)
	modPublic.GET("/authorization/callback", mod.V1OAuthCallback)
	modPublic.GET("/advertisings", mod.Advertisings)

	modPrivate := private.Group("/mod/:serverId/v1").Use(middleware.ServerMiddleware())
	modPrivate.GET("/authorization/check", mod.V1CheckAuthorization)
	modPrivate.POST("/authorization/logout-all", mod.V1LogoutAll)
	modPrivate.POST("/calling", mod.V1Calling)

	logger.Info("HTTP сервер запущен", "addr", ":4042")
	if err := engine.Run(":4042"); err != nil {
		logger.Fatal("HTTP сервер остановлен с ошибкой", "err", err)
	}
}
