package security

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
  signInKey = "hjdsfhsjd12&*"
  tokenTTL = 20 * time.Second
)

func GenerateHashPassword(password string) (string, error) {
  passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", err
  }
  return string(passwordHash), nil
}

func CheckPassword(passwordHash string, password string) error {
  err:=bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
  if err!=nil {
    return err
  }
  return nil
}

func SignIn (userName string, password string) error {
  var user models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return errors.New("incorrect username or password")
}
  if err := CheckPassword(user.PasswordHash, password); err!= nil {
    return errors.New("incorrect username or password")
  }

  return nil
}

func TokenFromHTTPRequest(r *http.Request) string {
  reqToken := r.Header.Get("Authorization")
  var tokenString string

  splitToken := strings.Split(reqToken, "Bearer ")
  if len(splitToken)>1 {
    tokenString = splitToken[1]
  }
  return tokenString
}

func GenerateToken (userName string, ttl time.Duration) (*Token, error) {
  var user models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return nil, errors.New("incorrect username")
}
  claims, err := NewClaims(user.ID, ttl)
  if err !=nil {
    return nil, err
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
    jwt.RegisteredClaims {
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
    IssuedAt: jwt.NewNumericDate(time.Now()),
    },
    user.ID})

    signedToken, errSign := token.SignedString([]byte(signInKey))
    if errSign != nil {
      return nil, errSign
    }

  return &Token{Value: signedToken, Claims: claims} , nil
}

func ParseToken(accessToken string) (uint, error) {
  token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("invalid signing method")
    }

    return []byte(signInKey), nil
  })
  if err != nil {
    return 0, err
  }
  
  claims, ok := token.Claims.(*tokenClaims)
  if !ok {
    return 0, errors.New("token claims are not of type")
  }

  return claims.UserID, nil
}