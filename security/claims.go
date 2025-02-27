package security

import (
	"myfavouritemovies/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewClaims(userID uint, duration time.Duration) (*models.TokenClaims, error) {
  tokenID, err := uuid.NewRandom()
  if err != nil {
    return nil, err
  }
  return &models.TokenClaims{
      RegisteredClaims: jwt.RegisteredClaims{
        ID:        tokenID.String(),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
      },
      UserID: userID,
    },
    nil
}
