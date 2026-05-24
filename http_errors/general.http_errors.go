// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"

	"sptools-backend/models"
)

var (
	FunctionDisabled  = models.NewCustomError(http.StatusForbidden, "This function is disabled for this server", "error.client.function.disabled")
	InvalidBody       = models.NewCustomError(http.StatusBadRequest, "Invalid body", "error.client.body")
)
