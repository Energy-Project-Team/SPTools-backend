// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package config

import (
	spwLib "github.com/xligenda/spworlds"
)

func GetSPworldsClientByID(id string) (*spwLib.Client, bool) {
	serverConfig, exist := GetServerByID(id)
	if !exist {
		return nil, false
	}
	return spwLib.NewClient(serverConfig.Card.ID, serverConfig.Card.Token), true
}
