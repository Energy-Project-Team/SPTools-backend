// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package models

import (
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	Code       string `json:"code"`
	StatusCode int    `json:"statusCode"`
}

type CustomErrorResponse struct {
	Message    string `json:"message"`
	Code       string `json:"code"`
	StatusCode int    `json:"statusCode"`
}

func (e *CustomErrorResponse) Error() string {
	return http.StatusText(e.StatusCode)
}

func NewCustomError(statusCode int, message, code string) *CustomErrorResponse {
	return &CustomErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}

type CustomResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func NewCustomResponse(statusCode int, message string) CustomResponse {
	return CustomResponse{StatusCode: statusCode, Message: message}
}

// /api/info
type InfoResponse struct {
	Node            Node    `json:"node"`
	AppVersion      string  `json:"version"`
	Service         string  `json:"service"`
	DiscordClientID string  `json:"discordClientID"`
	StartAt         string  `json:"startAt"`
	Uptime          float64 `json:"uptime"`
}

// /api/:version/version
type ReleaseResponse struct {
	Name          string `json:"name"`
	VersionNumber string `json:"versionNumber"`
	DatePublished string `json:"datePublished"`
	Files         []File `json:"files"`
	Changelog     string `json:"changelog"`
}

type File struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}
