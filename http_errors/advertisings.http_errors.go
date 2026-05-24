// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package http_errors

import (
	"net/http"

	"sptools-backend/models"
)

var (
	AdvertisingsGetData        = models.NewCustomError(http.StatusInternalServerError, "Couldn't get the data", "error.server.advertisings.get")
	AdvertisingsDataProcessing = models.NewCustomError(http.StatusInternalServerError, "Data processing error", "error.server.advertisings.processing")
)
