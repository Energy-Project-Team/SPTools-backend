// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package models

import "time"

type Node string

const (
	Production  Node = "production"
	Development Node = "development"
)

type Config struct {
	AppVersion  string
	Node        Node
	GithubToken string

	JWTSecretKey          string
	JWTTokenExpireMinutes int

	MongoInitDBHost         string
	MongoInitDBRootUsername string
	MongoInitDBRootPassword string

	Servers []ServerConfig

	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURL  string
	DiscordAPIOrigin    string

	StartAt time.Time
}

type ServerConfig struct {
	ID          string               `yaml:"id" json:"id"`
	Name        string               `yaml:"name" json:"name"`
	DisplayName string               `yaml:"displayName" json:"displayName"`
	IPs         []string             `yaml:"ips" json:"ips"`
	Enabled     bool                 `yaml:"enabled" json:"enabled"`
	Card        ServerCardConfig     `yaml:"card" json:"-"`
	Features    ServerFeaturesConfig `yaml:"features" json:"features"`
	Calling     ServerCallingConfig  `yaml:"calling" json:"calling"`
}

type ServerCardConfig struct {
	ID    string `yaml:"id" json:"-"`
	Token string `yaml:"token" json:"-"`
}

type ServerFeaturesConfig struct {
	Calling      bool `yaml:"calling" json:"calling"`
	Advertisings bool `yaml:"advertisings" json:"advertisings"`
}

type ServerCallingConfig struct {
	Webhooks []ServerCallingWebhookConfig `yaml:"webhooks" json:"webhooks"`
}

type ServerCallingWebhookConfig struct {
	Service    string `yaml:"service" json:"service"`
	Name       string `yaml:"name" json:"name"`
	RoleID     string `yaml:"roleId" json:"-"`
	Webhook    string `yaml:"webhook" json:"-"`
	ImageURL   string `yaml:"imageUrl" json:"imageUrl"`
}
