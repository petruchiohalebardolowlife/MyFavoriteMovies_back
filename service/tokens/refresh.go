package service

import (
	"errors"
	"myfavouritemovies/repository"
)

func Refresh(refreshToken string) (string, string, error) {
  claims, err := Validate(refreshToken)
  if err != nil {
    return "", "", err
  }

  if err := repository.CheckTokenInBlackList(claims.ID); err != nil {
    return "", "", err
  }

  newAccessToken, newRefreshToken, newError := Generate(claims.UserID)
  if newError != nil {
    return "", "", errors.New("failed to generate tokens")
  }

  return newAccessToken, newRefreshToken, nil
}
