package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenClaims struct {
  jwt.RegisteredClaims
  UserID uint `json:"user_id"`
}

type Token struct {
  Value string
  Claims *tokenClaims
}


func NewClaims(id uint, duration time.Duration) (*tokenClaims, error) {
  tokenID, err := uuid.NewRandom()
  if err != nil {
    return nil, err
  }
  return &tokenClaims{
    RegisteredClaims : jwt.RegisteredClaims{
      ID: tokenID.String(),
      IssuedAt: jwt.NewNumericDate(time.Now()),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
    },
    UserID: id,
  },
  nil
}