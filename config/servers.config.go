// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package config

import (
	"strings"

	"sptools-backend/models"
)

func GetEnabledServers() []models.ServerConfig {
	servers := make([]models.ServerConfig, 0)

	for _, server := range GlobalConfig.Servers {
		if server.Enabled {
			servers = append(servers, server)
		}
	}

	return servers
}

func GetServerByID(id string) (*models.ServerConfig, bool) {
	for index := range GlobalConfig.Servers {
		server := &GlobalConfig.Servers[index]

		if !server.Enabled {
			continue
		}

		if strings.EqualFold(server.ID, id) {
			return server, true
		}
	}

	return nil, false
}

func HasServer(id string) bool {
	_, ok := GetServerByID(id)
	return ok
}

func ServerHasCalling(id string) bool {
	server, ok := GetServerByID(id)
	if !ok {
		return false
	}

	return server.Features.Calling
}

func ServerHasAdvertisings(id string) bool {
	server, ok := GetServerByID(id)
	if !ok {
		return false
	}

	return server.Features.Advertisings
}

func GetServerCallingWebhook(serverID string, service string) (*models.ServerCallingWebhookConfig, bool) {
	server, ok := GetServerByID(serverID)
	if !ok {
		return nil, false
	}

	if !server.Features.Calling {
		return nil, false
	}

	for index := range server.Calling.Webhooks {
		webhook := &server.Calling.Webhooks[index]

		if strings.EqualFold(webhook.Service, service) {
			return webhook, true
		}
	}

	return nil, false
}

func GetServerCallingWebhooks(serverID string) []models.ServerCallingWebhookConfig {
	server, ok := GetServerByID(serverID)
	if !ok || !server.Features.Calling {
		return []models.ServerCallingWebhookConfig{}
	}

	return server.Calling.Webhooks
}
