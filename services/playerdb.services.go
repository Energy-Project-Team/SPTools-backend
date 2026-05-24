// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only

package services

import (
	"fmt"
	"net/http"
	"time"

	"sptools-backend/models"

	"github.com/imroc/req/v3"
)

const playerDBMinecraftProfileURL = "https://playerdb.co/api/player/minecraft/%s"

type PlayerDBService struct {
	client *req.Client
}

func NewPlayerDBService() *PlayerDBService {
	return &PlayerDBService{
		client: req.C().
			SetTimeout(15 * time.Second).
			SetCommonHeader("User-Agent", "Energy-Project-Team/SPTools-backend"),
	}
}

func (s *PlayerDBService) GetProfileByNickname(nickname string) (models.PlayerDBAPIPlayer, error) {
	return s.getMinecraftProfile(nickname)
}

func (s *PlayerDBService) GetProfileByUUID(uuid string) (models.PlayerDBAPIPlayer, error) {
	return s.getMinecraftProfile(uuid)
}

func (s *PlayerDBService) getMinecraftProfile(query string) (models.PlayerDBAPIPlayer, error) {
	var apiResp models.PlayerDBAPIResponse

	resp, err := s.client.R().
		SetSuccessResult(&apiResp).
		Get(fmt.Sprintf(playerDBMinecraftProfileURL, query))

	if err != nil {
		return models.PlayerDBAPIPlayer{}, fmt.Errorf("ошибка запроса к PlayerDB API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.PlayerDBAPIPlayer{}, fmt.Errorf(
			"PlayerDB API вернул статус %d: %s",
			resp.StatusCode,
			resp.String(),
		)
	}

	if !apiResp.Success {
		return models.PlayerDBAPIPlayer{}, fmt.Errorf(
			"PlayerDB API вернул неуспешный ответ: %s",
			apiResp.Message,
		)
	}

	return apiResp.Data.Player, nil
}