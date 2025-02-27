package security

import (
	"errors"
	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
  passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", err
  }

  return string(passwordHash), nil
}

func CheckPassword(passwordHash string, password string) error {
  err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
  if err != nil {
    return err
  }

  return nil
}

func SignIn(userName string, password string) error {
  var user models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return err
  }
  if err := CheckPassword(user.PasswordHash, password); err != nil {
    return err
  }

  return nil
}

func TokenFromCookie(r *http.Request) (string, error) {
  reqToken, err := r.Cookie("jwtRefresh")
  if err != nil {
    return "", err
  }

  return reqToken.Value, nil
}

func GenerateToken(userID uint, ttl time.Duration) (*models.Token, error) {
  claims, err := NewClaims(userID, ttl)
  if err != nil {
    return nil, err
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  signedToken, errSign := token.SignedString([]byte(config.TOKEN_KEY))
  if errSign != nil {
    return nil, errSign
  }

  return &models.Token{Value: signedToken, Claims: claims}, nil
}

func ParseToken(TokenValue string) (*models.TokenClaims, error) {
  token, err := jwt.ParseWithClaims(TokenValue, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("invalid signing method")
    }

    return []byte(config.TOKEN_KEY), nil
  })
  if err != nil {
    return nil, err
  }

  claims, ok := token.Claims.(*models.TokenClaims)
  if !ok {
    return nil, errors.New("token claims are not of type")
  }

  return claims, nil
}

func UpdateTokens(userID uint, ttlAccess, ttlRefresh time.Duration) (*models.Tokens, error) {
  newAccessToken, errAccessToken := GenerateToken(userID, ttlAccess)
  if errAccessToken != nil {
    return nil, errAccessToken
  }
  newRefreshToken, errRefreshToken := GenerateToken(userID, ttlRefresh)
  if errRefreshToken != nil {
    return nil, errRefreshToken
  }
  return &models.Tokens{Access: newAccessToken, Refresh: newRefreshToken}, nil
}

func SetTokenInCookie(writer http.ResponseWriter, refreshToken string) {
  http.SetCookie(writer, &http.Cookie{
    Name:     "jwtRefresh",
    Value:    refreshToken,
    Path:     "/",
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
  })
}

func DeleteTokensFromCookie(writer http.ResponseWriter) {
  http.SetCookie(writer, &http.Cookie{
    Name:     "jwt_access_token",
    Path:     "/",
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   -1,
  })

  http.SetCookie(writer, &http.Cookie{
    Name:     "jwt_refresh_token",
    Path:     "/",
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   -1,
  })
}
