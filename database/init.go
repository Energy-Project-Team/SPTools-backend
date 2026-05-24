// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package database

import (
	"context"

	"sptools-backend/logger"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
	Engine  *gin.Engine
	MongoDB *mongo.Database
}

func SetupDataBase(engine *gin.Engine, mongoDB *mongo.Database) *DataBase {
	return &DataBase{Engine: engine, MongoDB: mongoDB}
}

func (db *DataBase) Collection(name string) *mongo.Collection {
	return db.MongoDB.Collection(name)
}

func (db *DataBase) CreateIndexes(ctx context.Context) error {
	collection := db.MongoDB.Collection("Users")
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "minecraftUUID", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "discordId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "servers.id", Value: 1}}},
	})
	if err != nil {
		logger.Error("Не удалось создать индексы в коллекции", "collectionName", collection.Name(), "err", err)
		return err
	}
	return nil
}
