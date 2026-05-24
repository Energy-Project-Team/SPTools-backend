// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package services

import (
	"fmt"

	"sptools-backend/config"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) CreateToken(claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный способ подписи: %v", token.Header["alg"])
		}
		return []byte(config.GlobalConfig.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("недействительный токен")
}
