// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package database

import (
	"context"
	"strings"
	"time"

	"sptools-backend/http_errors"
	"sptools-backend/logger"
	"sptools-backend/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DataBase) usersCollection() *mongo.Collection {
	return db.MongoDB.Collection("Users")
}

func normalizeMinecraftUUID(minecraftUUID string) string {
	return strings.ReplaceAll(strings.TrimSpace(minecraftUUID), "-", "")
}

func (db *DataBase) Users(ctx context.Context) ([]models.User, error) {
	cursor, err := db.usersCollection().Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		logger.Error("Не удалось создать запрос на получение пользователей", "err", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	documents := make([]models.User, 0)
	if err := cursor.All(ctx, &documents); err != nil {
		logger.Error("Не удалось получить пользователей", "err", err)
		return nil, err
	}
	return documents, nil
}

func (db *DataBase) UserByID(ctx context.Context, id string) (*models.User, error) {
	var document models.User
	err := db.usersCollection().FindOne(ctx, bson.M{"id": strings.TrimSpace(id)}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http_errors.UserNotFound
		}
		logger.Error("Не удалось получить пользователя по id", "err", err, "id", id)
		return nil, err
	}
	return &document, nil
}

func (db *DataBase) UserByMinecraftUUID(ctx context.Context, minecraftUUID string) (*models.User, error) {
	var document models.User
	err := db.usersCollection().FindOne(ctx, bson.M{"minecraftUUID": normalizeMinecraftUUID(minecraftUUID)}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http_errors.UserNotFound
		}
		logger.Error("Не удалось получить пользователя по minecraftUUID", "err", err, "minecraftUUID", minecraftUUID)
		return nil, err
	}
	return &document, nil
}

func (db *DataBase) UserByDiscordID(ctx context.Context, discordID string) (*models.User, error) {
	var document models.User
	err := db.usersCollection().FindOne(ctx, bson.M{"discordId": strings.TrimSpace(discordID)}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http_errors.UserNotFound
		}
		logger.Error("Не удалось получить пользователя по discordId", "err", err, "discordId", discordID)
		return nil, err
	}
	return &document, nil
}

func (db *DataBase) UsersCreate(ctx context.Context, user models.User) error {
	now := time.Now()
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	if user.Status == "" {
		user.Status = "active"
	}
	if len(user.Roles) == 0 {
		user.Roles = []string{"user"}
	}
	if user.CreatedAt == nil || user.CreatedAt.IsZero() {
		user.CreatedAt = &now
	}
	if user.UpdatedAt == nil || user.UpdatedAt.IsZero() {
		user.UpdatedAt = &now
	}
	_, err := db.usersCollection().InsertOne(ctx, user)
	if err != nil {
		logger.Error("Не удалось создать пользователя", "err", err, "user", user)
		return err
	}
	return nil
}

func (db *DataBase) UserUpdate(ctx context.Context, user models.User) error {
	now := time.Now()
	user.UpdatedAt = &now
	result, err := db.usersCollection().UpdateOne(ctx, bson.M{"id": user.ID}, bson.M{"$set": user})
	if err != nil {
		logger.Error("Не удалось обновить пользователя", "err", err, "user", user)
		return err
	}
	if result.MatchedCount == 0 {
		return http_errors.UserNotFound
	}
	return nil
}

func (db *DataBase) UpsertUserLogin(ctx context.Context, minecraftUUID string, discordID string, serverID string) (*models.User, error) {
	now := time.Now()
	minecraftUUID = normalizeMinecraftUUID(minecraftUUID)
	discordID = strings.TrimSpace(discordID)
	serverID = strings.TrimSpace(serverID)

	user, err := db.UserByMinecraftUUID(ctx, minecraftUUID)
	if err != nil && err != http_errors.UserNotFound {
		return nil, err
	}

	if err == http_errors.UserNotFound {
		user = &models.User{ID: uuid.New().String(), MinecraftUUID: minecraftUUID, DiscordID: discordID, Servers: []models.UserServer{{ID: serverID, FirstLoginAt: &now, LastLoginAt: &now}}, Status: "active", Roles: []string{"user"}, TokenVersion: 0, CreatedAt: &now, UpdatedAt: &now, LastLoginAt: &now}
		if err := db.UsersCreate(ctx, *user); err != nil {
			return nil, err
		}
		return user, nil
	}

	servers := user.Servers
	serverFound := false
	for index := range servers {
		if servers[index].ID == serverID {
			serverFound = true
			servers[index].LastLoginAt = &now
			if servers[index].FirstLoginAt == nil {
				servers[index].FirstLoginAt = &now
			}
			break
		}
	}
	if !serverFound {
		servers = append(servers, models.UserServer{ID: serverID, FirstLoginAt: &now, LastLoginAt: &now})
	}

	update := bson.M{"$set": bson.M{"discordId": discordID, "servers": servers, "updatedAt": now, "lastLoginAt": now}}
	result, err := db.usersCollection().UpdateOne(ctx, bson.M{"id": user.ID}, update)
	if err != nil {
		logger.Error("Не удалось обновить пользователя после авторизации", "err", err, "userId", user.ID, "serverId", serverID)
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, http_errors.UserNotFound
	}
	user.DiscordID = discordID
	user.Servers = servers
	user.UpdatedAt = &now
	user.LastLoginAt = &now
	return user, nil
}

func (db *DataBase) IncrementUserTokenVersion(ctx context.Context, userID string) error {
	now := time.Now()
	result, err := db.usersCollection().UpdateOne(ctx, bson.M{"id": strings.TrimSpace(userID)}, bson.M{"$inc": bson.M{"tokenVersion": int64(1)}, "$set": bson.M{"updatedAt": now}})
	if err != nil {
		logger.Error("Не удалось увеличить tokenVersion пользователя", "err", err, "userID", userID)
		return err
	}
	if result.MatchedCount == 0 {
		return http_errors.UserNotFound
	}
	return nil
}

func (db *DataBase) UserDelete(ctx context.Context, id string) error {
	result, err := db.usersCollection().DeleteOne(ctx, bson.M{"id": strings.TrimSpace(id)})
	if err != nil {
		logger.Error("Не удалось удалить пользователя по id", "err", err, "id", id)
		return err
	}
	if result.DeletedCount == 0 {
		return http_errors.UserNotFound
	}
	return nil
}

func (db *DataBase) UserDeleteByMinecraftUUID(ctx context.Context, minecraftUUID string) error {
	result, err := db.usersCollection().DeleteOne(ctx, bson.M{"minecraftUUID": normalizeMinecraftUUID(minecraftUUID)})
	if err != nil {
		logger.Error("Не удалось удалить пользователя по minecraftUUID", "err", err, "minecraftUUID", minecraftUUID)
		return err
	}
	if result.DeletedCount == 0 {
		return http_errors.UserNotFound
	}
	return nil
}
