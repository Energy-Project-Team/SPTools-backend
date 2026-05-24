// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package models

type CallingSuccessResponse struct {
	Message       string `json:"message"`
	Service       string `json:"service"`
	DiscordID     string `json:"discordID"`
	Nickname      string `json:"nickname"`
	MinecraftUUID string `json:"minecraftUUID"`
}

type CallingFailureResponse struct {
	Message         string                 `json:"message"`
	Service         string                 `json:"service"`
	DiscordID       string                 `json:"discordID"`
	Nickname        string                 `json:"nickname"`
	MinecraftUUID   string                 `json:"minecraftUUID"`
	DiscordResponse map[string]interface{} `json:"discordResponse"`
}
