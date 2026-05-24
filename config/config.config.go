// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sptools-backend/logger"
	"sptools-backend/models"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var GlobalConfig models.Config

type yamlConfig struct {
	Servers []models.ServerConfig `yaml:"servers"`
}

func LoadConfig(appVersion string) error {
	if err := godotenv.Overload(); err != nil {
		logger.Warn("Файл .env не найден")
	} else {
		logger.Info("Файл .env найден")
	}

	logger.Info("Проверяю переменные окружения")

	node := models.Node(os.Getenv("NODE"))

	GlobalConfig.AppVersion = appVersion
	if strings.ToLower(string(node)) == string(models.Production) {
		GlobalConfig.Node = models.Production
	} else {
		GlobalConfig.Node = models.Development
	}

	requiredKeys := []string{
		"NODE",
		"CONFIG_PATH",

		"JWT_SECRET_KEY",
		"JWT_TOKEN_EXPIRE_MINUTES",

		"GITHUB_TOKEN",

		"MONGO_INITDB_HOST",
		"MONGO_INITDB_ROOT_USERNAME",
		"MONGO_INITDB_ROOT_PASSWORD",

		"DISCORD_CLIENT_ID",
		"DISCORD_CLIENT_SECRET",
	}

	var missingKeys []string
	for _, key := range requiredKeys {
		if os.Getenv(key) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		logger.Fatal(fmt.Sprintf("Не найдены следующие ключи в переменных окружения: %s", strings.Join(missingKeys, ", ")))
	}

	GlobalConfig.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

	tokenTTL, err := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRE_MINUTES"))
	if err != nil {
		tokenTTL = 43200
	}

	GlobalConfig.JWTTokenExpireMinutes = tokenTTL

	GlobalConfig.MongoInitDBHost = os.Getenv("MONGO_INITDB_HOST")
	GlobalConfig.MongoInitDBRootUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	GlobalConfig.MongoInitDBRootPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	GlobalConfig.DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
	GlobalConfig.DiscordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")

	if GlobalConfig.Node == models.Development {
		GlobalConfig.DiscordRedirectURL = "http://localhost:4011/api/mod/%f/v1/authorization/callback"
	} else {
		GlobalConfig.DiscordRedirectURL = "https://sptools.energyproject.dev/api/mod/%f/v1/authorization/callback"
	}

	GlobalConfig.GithubToken = os.Getenv("GITHUB_TOKEN")

	if err := loadYAMLConfig(os.Getenv("CONFIG_PATH")); err != nil {
		return err
	}

	validateServers()

	GlobalConfig.StartAt = time.Now()

	logger.Info("Конфигурация успешно загружена")

	return nil
}

func loadYAMLConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("не удалось прочитать YAML-конфиг %s: %w", path, err)
	}

	var cfg yamlConfig
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return fmt.Errorf("не удалось распарсить YAML-конфиг %s: %w", path, err)
	}

	GlobalConfig.Servers = cfg.Servers

	logger.Info(fmt.Sprintf("Загружено серверов из YAML-конфига: %d", len(GlobalConfig.Servers)))

	return nil
}

func validateServers() {
	if len(GlobalConfig.Servers) == 0 {
		logger.Fatal("В YAML-конфиге не указан ни один сервер")
	}

	ids := make(map[string]bool)

	for _, server := range GlobalConfig.Servers {
		if !server.Enabled {
			logger.Info(fmt.Sprintf("Сервер %s выключен в конфиге", server.ID))
			continue
		}

		if server.ID == "" {
			logger.Fatal("В YAML-конфиге есть включённый сервер без id")
		}

		if ids[server.ID] {
			logger.Fatal(fmt.Sprintf("Дублирующийся id сервера в YAML-конфиге: %s", server.ID))
		}

		ids[server.ID] = true

		if server.Name == "" {
			logger.Fatal(fmt.Sprintf("У сервера %s не указан name", server.ID))
		}

		if server.DisplayName == "" {
			logger.Fatal(fmt.Sprintf("У сервера %s не указан displayName", server.ID))
		}

		if len(server.IPs) == 0 {
			logger.Fatal(fmt.Sprintf("У сервера %s не указаны ips", server.ID))
		}

		if server.Card.ID == "" {
			logger.Fatal(fmt.Sprintf("У сервера %s не указан card.id", server.ID))
		}

		if server.Card.Token == "" {
			logger.Fatal(fmt.Sprintf("У сервера %s не указан card.token", server.ID))
		}

		if server.Features.Calling {
			if len(server.Calling.Webhooks) == 0 {
				logger.Fatal(fmt.Sprintf("У сервера %s включён calling, но calling.webhooks пустой", server.ID))
			}

			services := make(map[string]bool)

			for _, webhook := range server.Calling.Webhooks {
				if webhook.Service == "" {
					logger.Fatal(fmt.Sprintf("У сервера %s есть calling webhook без service", server.ID))
				}

				if services[webhook.Service] {
					logger.Fatal(fmt.Sprintf("У сервера %s дублирующийся calling service: %s", server.ID, webhook.Service))
				}

				services[webhook.Service] = true

				if webhook.Name == "" {
					logger.Fatal(fmt.Sprintf("У сервера %s у calling service %s не указан name", server.ID, webhook.Service))
				}

				if webhook.RoleID == "" {
					logger.Fatal(fmt.Sprintf("У сервера %s у calling service %s не указан roleId", server.ID, webhook.Service))
				}

				if webhook.Webhook == "" {
					logger.Fatal(fmt.Sprintf("У сервера %s у calling service %s не указан webhook", server.ID, webhook.Service))
				}
			}
		}

		logger.Info(fmt.Sprintf("Сервер включён: %s (%s), calling=%t, advertisings=%t", server.ID, server.DisplayName, server.Features.Calling, server.Features.Advertisings))
	}
}
