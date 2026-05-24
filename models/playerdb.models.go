// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package models

type PlayerDBAPIProperty struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature,omitempty"`
}

type PlayerDBAPIMeta struct {
	CachedAt int64 `json:"cached_at"`
}

type PlayerDBAPIPlayer struct {
	Meta        PlayerDBAPIMeta       `json:"meta"`
	Username    string                `json:"username"`
	ID          string                `json:"id"`
	RawID       string                `json:"raw_id"`
	Avatar      string                `json:"avatar"`
	SkinTexture string                `json:"skin_texture"`
	Properties  []PlayerDBAPIProperty `json:"properties"`
	NameHistory []string              `json:"name_history"`
}

type PlayerDBAPIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Player PlayerDBAPIPlayer `json:"player"`
	} `json:"data"`
	Success bool `json:"success"`
}
