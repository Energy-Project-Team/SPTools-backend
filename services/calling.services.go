// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sptools-backend/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type CallingService struct {
	Client  *http.Client
	MongoDB *mongo.Database
}

func NewCallingService() *CallingService {
	return &CallingService{
		Client: &http.Client{},
	}
}

func (s *CallingService) Calling(webhookURL, nickname, discordID, roleID, comment, coordinates string) (map[string]interface{}, int, error) {
	sendInfo := fmt.Sprintf("Webhook URL: %s; nickname: %s; discordID: %s; roleID: %s; comment: %s; coordinates: %s",
		webhookURL, nickname, discordID, roleID, comment, coordinates,
	)

	var text string
	if coordinates != "" {
		text = fmt.Sprintf("`%s` (<@%s>) вызывает <@&%s> на: `%s` **%s**",
			nickname, discordID, roleID, comment, coordinates)
	} else {
		text = fmt.Sprintf("`%s` (<@%s>) вызывает <@&%s> на: `%s`",
			nickname, discordID, roleID, comment)
	}
	body, err := json.Marshal(map[string]interface{}{"content": text})
	if err != nil {
		return nil, 0, err
	}

	for {
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(body))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.Client.Do(req)
		if err != nil {
			return nil, 0, err
		}
		defer resp.Body.Close()

		var responseData map[string]interface{}
		if resp.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
				responseData = map[string]interface{}{}
			}
		}

		if resp.StatusCode == http.StatusNoContent {
			logger.Info("Успешно отправлен webhook", "details", sendInfo)
			return responseData, resp.StatusCode, nil
		} else if resp.StatusCode == 429 {
			retryAfter, ok := responseData["retry_after"].(float64)
			if ok {
				logger.Warn("превышен лимит отправки запросов к Discord, ожидание", "retryAfterMs", retryAfter, "details", sendInfo)
				time.Sleep(time.Duration(retryAfter) * time.Millisecond)
				continue
			} else {
				logger.Warn("превышен лимит отправки запросов к Discord, ожидание 10 секунд", "details", sendInfo)
				time.Sleep(10 * time.Second)
				continue
			}
		} else {
			logger.Error("неизвестный ответ от Discord", "responseData", responseData, "status", resp.StatusCode, "details", sendInfo)
			return responseData, resp.StatusCode, nil
		}
	}
}
