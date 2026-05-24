// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package general_controllers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"sptools-backend/config"
	"sptools-backend/http_errors"
	"sptools-backend/logger"
	"sptools-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

const (
	modrinthVersionsURL = "https://api.modrinth.com/v2/project/spmhelper/version"
	// modrinthVersionsURL = "https://api.modrinth.com/v2/project/sptools/version"
	versionsCacheTTL = 60 * time.Second
)

var versionsCache = struct {
	sync.RWMutex
	data      []models.ReleaseResponse
	timestamp time.Time
}{}

func (ic *GeneralControllers) V1Versions(c *gin.Context) {
	releases, err := ic.getModrinthVersions()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, releases)
}

func (ic *GeneralControllers) V1LatestVersion(c *gin.Context) {
	releases, err := ic.getModrinthVersions()
	if err != nil {
		c.Error(err)
		return
	}

	if len(releases) == 0 {
		logger.Warn("версии не найдены")
		c.Error(http_errors.VersionsNotFound)
		return
	}

	c.JSON(http.StatusOK, releases[0])
}

func (ic *GeneralControllers) getModrinthVersions() ([]models.ReleaseResponse, error) {
	currentTime := time.Now()

	versionsCache.RLock()
	cachedData := versionsCache.data
	cachedTimestamp := versionsCache.timestamp
	versionsCache.RUnlock()

	if len(cachedData) > 0 && currentTime.Sub(cachedTimestamp) < versionsCacheTTL {
		return cachedData, nil
	}

	var releases []models.ReleaseResponse

	resp, err := req.C().
		SetTimeout(15*time.Second).
		R().
		SetHeader("Authorization", config.GlobalConfig.GithubToken).
		SetHeader("Content-Type", "application/json").
		SetHeader(
			"User-Agent",
			fmt.Sprintf(
				"Energy-Project-Team/SPwHelper-backend/%s (sptools.energyproject.dev)",
				config.GlobalConfig.AppVersion,
			),
		).
		SetSuccessResult(&releases).
		Get(modrinthVersionsURL)

	if err != nil {
		logger.Error("не удалось выполнить запрос в Modrinth API", "err", err)
		return nil, http_errors.VersionsSendModrinthAPI
	}

	if !resp.IsSuccessState() {
		logger.Error("ошибка поиска версии", "statusCode", resp.StatusCode, "body", resp.String())
		return nil, http_errors.VersionsSearchVer
	}

	versionsCache.Lock()
	versionsCache.data = releases
	versionsCache.timestamp = currentTime
	versionsCache.Unlock()

	return releases, nil
}
