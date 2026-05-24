// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package models

import "time"

type User struct {
	ID             string       `bson:"id" json:"id"`
	MinecraftUUID  string       `bson:"minecraftUUID" json:"minecraftUUID"`
	DiscordID      string       `bson:"discordId" json:"discordId"`
	Servers        []UserServer `bson:"servers" json:"servers"`
	Status         string       `bson:"status" json:"status"`
	Roles          []string     `bson:"roles" json:"roles"`
	TokenVersion   int64        `bson:"tokenVersion" json:"-"`
	IsBanned       bool         `bson:"isBanned" json:"isBanned,omitempty"`
	BanReason      *string      `bson:"banReason" json:"banReason,omitempty"`
	BannedAt       *time.Time   `bson:"bannedAt" json:"bannedAt,omitempty"`
	BannedByUserID *string      `bson:"bannedByUserID" json:"bannedByUserID,omitempty"`
	CreatedAt      *time.Time   `bson:"createdAt" json:"-"`
	UpdatedAt      *time.Time   `bson:"updatedAt" json:"-"`
	LastLoginAt    *time.Time   `bson:"lastLoginAt" json:"-"`
}

type UserServer struct {
	ID           string     `bson:"id" json:"id"`
	FirstLoginAt *time.Time `bson:"firstLoginAt" json:"firstLoginAt"`
	LastLoginAt  *time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
}

func (user *User) HasServer(serverID string) bool {
	for _, server := range user.Servers {
		if server.ID == serverID {
			return true
		}
	}
	return false
}

func (user *User) ClearAdminInformation() {
	user.TokenVersion = 0
	user.IsBanned = false
	user.BanReason = nil
	user.BannedAt = nil
	user.BannedByUserID = nil
	user.CreatedAt = nil
	user.UpdatedAt = nil
}
