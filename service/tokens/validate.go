package service

import (
	"errors"
	config "myfavouritemovies/configs"
	"myfavouritemovies/models"

	"github.com/golang-jwt/jwt/v5"
)

func Validate(inputToken string) (*models.TokenClaims, error) {
  token, err := jwt.ParseWithClaims(inputToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("invalid signing method")
    }

    return []byte(config.TOKEN_KEY), nil
  })
  if err != nil {
    return nil, errors.New("token expired")
  }

  claims, ok := token.Claims.(*models.TokenClaims)
  if !ok {
    return nil, errors.New("token claims are not of type")
  }

  return claims, nil
}
