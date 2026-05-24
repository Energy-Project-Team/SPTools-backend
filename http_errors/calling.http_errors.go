// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"

	"sptools-backend/models"
)

var (
	CallingGetNickname    = models.NewCustomError(http.StatusInternalServerError, "Couldn't get the player's nickname", "error.server.calling.nickname")
	CallingSendWebhook    = models.NewCustomError(http.StatusInternalServerError, "Error sending webhook", "error.server.calling.webhook")
)
