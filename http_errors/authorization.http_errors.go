// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"

	"sptools-backend/models"
)

var (
	AuthorizationInvalidToken     = models.NewCustomError(http.StatusUnauthorized, "Invalid token", "error.client.authorization.token.invalid")
	AuthorizationCallbackCode     = models.NewCustomError(http.StatusBadRequest, "Missing authorization code", "error.client.authorization.code.required")
	AuthorizationCallbackState    = models.NewCustomError(http.StatusBadRequest, "Missing state", "error.client.authorization.state.required")
	AuthorizationInvalidState     = models.NewCustomError(http.StatusBadRequest, "Invalid state", "error.client.authorization.state.invalid")
	AuthorizationSSENotFound      = models.NewCustomError(http.StatusBadRequest, "SSE connection not found", "error.client.authorization.session.notFound")
	AuthorizationAccessDenied     = models.NewCustomError(http.StatusUnauthorized, "Authorization denied by user", "error.client.authorization.denied")
	AuthorizationTokenRequest     = models.NewCustomError(http.StatusInternalServerError, "Failed to obtain Discord token", "error.server.authorization.discordToken")
	AuthorizationUserDataRequest  = models.NewCustomError(http.StatusInternalServerError, "Failed to obtain Discord user data", "error.server.authorization.discordUser")
	AuthorizationCreateToken      = models.NewCustomError(http.StatusInternalServerError, "Failed to create API token", "error.server.authorization.token.create")
	AuthorizationCreateUser       = models.NewCustomError(http.StatusInternalServerError, "Failed to create user", "error.server.authorization.user.create")
	AuthorizationFindUser         = models.NewCustomError(http.StatusInternalServerError, "Failed to find user", "error.server.authorization.user.find")
	AuthorizationUserNotFound     = models.NewCustomError(http.StatusForbidden, "User not found", "error.client.authorization.user.notFound")
	AuthorizationSPWorldsMismatch = models.NewCustomError(http.StatusForbidden, "Player not found on SP Worlds", "error.client.authorization.spworlds.mismatch")
	AuthorizationMissingToken     = models.NewCustomError(http.StatusUnauthorized, "Missing authorization token", "error.client.authorization.token.missing")
)
