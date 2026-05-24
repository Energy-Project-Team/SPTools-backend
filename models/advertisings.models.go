// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package models

type AdvertisingItem struct {
	Text     string  `bson:"text" json:"text"`
	ImageURL *string `bson:"imageURL" json:"imageURL"`
	URL      *string `bson:"url" json:"url"`
}
