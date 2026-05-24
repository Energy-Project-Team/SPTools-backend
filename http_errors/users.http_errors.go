// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"

	"sptools-backend/models"
)

var (
	UserNotFound = models.NewCustomError(http.StatusNotFound, "User not found", "error.client.user.not_found")
)
