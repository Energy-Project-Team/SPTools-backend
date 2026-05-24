// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"
	"sptools-backend/models"
)

var (
	ServerNotFound = models.NewCustomError(http.StatusNotFound, "Server not found", "error.client.server.not_found")
)
