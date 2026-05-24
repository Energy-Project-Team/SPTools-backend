// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package database

import (
	"context"
	"strings"

	"sptools-backend/logger"
	"sptools-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DataBase) advertisingCollection(serverID string) *mongo.Collection {
	return db.MongoDB.Collection(strings.TrimSpace(serverID) + "Advertising")
}

func (db *DataBase) Advertisings(ctx context.Context, serverID string) ([]models.AdvertisingItem, error) {
	cursor, err := db.advertisingCollection(serverID).Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		logger.Error("Не удалось создать запрос на получение рекламных объявлений", "err", err, "serverId", serverID)
		return nil, err
	}
	defer cursor.Close(ctx)

	documents := make([]models.AdvertisingItem, 0)
	if err := cursor.All(ctx, &documents); err != nil {
		logger.Error("Не удалось получить рекламные объявления", "err", err, "serverId", serverID)
		return nil, err
	}
	return documents, nil
}
