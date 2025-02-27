package service

import (
	"errors"
	config "myfavouritemovies/configs"
	"myfavouritemovies/models"
	"myfavouritemovies/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userID uint) (string, string, error) {
  claimsAccess, err := security.NewClaims(userID, 1*time.Minute)
  if err != nil {
    return "", "", err
  }
  accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
  signedAccessToken, errSign := accessToken.SignedString([]byte(config.TOKEN_KEY))
  if errSign != nil {
    return "", "", errSign
  }

  claimsRefresh, err := security.NewClaims(userID, 60*24*time.Hour)
  if err != nil {
    return "", "", err
  }
  refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
  signedRefreshToken, errSign := refreshToken.SignedString([]byte(config.TOKEN_KEY))
  if errSign != nil {
    return "", "", errSign
  }

  return signedAccessToken, signedRefreshToken, nil
}

func Validate(inputToken string) (uint, error) {
  token, err := jwt.ParseWithClaims(inputToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("invalid signing method")
    }

    return []byte(config.TOKEN_KEY), nil
  })
  if err != nil {
    return 0, errors.New("token expired")
  }

  claims, ok := token.Claims.(*models.TokenClaims)
  if !ok {
    return 0, errors.New("token claims are not of type")
  }

  return claims.UserID, nil
}

func Refresh(refreshToken string) (string, string, error) {
  userID, err := Validate(refreshToken)
  if err != nil {
    return "", "", err
  }

  newAccessToken, newRefreshToken, newError := GenerateTokens(userID)
  if newError != nil {
    return "", "", errors.New("failed to generate tokens")
  }

  return newAccessToken, newRefreshToken, nil
}
