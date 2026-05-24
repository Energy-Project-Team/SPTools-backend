// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"

	"sptools-backend/models"
)

var (
	VersionsSendModrinthAPI = models.NewCustomError(http.StatusBadGateway, "Couldn't complete the request in the Modrinth API", "error.server.versions.modrinth")
	VersionsSearchVer       = models.NewCustomError(http.StatusBadGateway, "Version search error", "error.server.versions.search")
	VersionsNotFound        = models.NewCustomError(http.StatusNotFound, "Versions not found", "error.client.versions.notFound")
)
