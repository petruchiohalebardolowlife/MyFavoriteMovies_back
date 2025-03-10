package service

import (
	config "myfavouritemovies/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(userID uint) (string, string, error) {
  claimsAccess, err := NewClaims(userID, 15*time.Minute)
  if err != nil {
    return "", "", err
  }
  accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
  signedAccessToken, errSign := accessToken.SignedString([]byte(config.TOKEN_KEY))
  if errSign != nil {
    return "", "", errSign
  }

  claimsRefresh, err := NewClaims(userID, 60*24*time.Hour)
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
