package security

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"net/http"
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

func FingerPrintFromHTTPRequest(r *http.Request) (string, error) {
  fingerprint := r.Header.Get("Fingerprint")
  if len(fingerprint) == 0 {
    return "", errors.New("no fingerprint")
  }
  return fingerprint, nil
}



func TokenFromCookie(r *http.Request, tokentype string) string {
  reqToken, err := r.Cookie(tokentype)
  if err != nil {
    return ""
  }
  
  return reqToken.Value
}

func GenerateToken (userID uint, ttl time.Duration) (*Token, error) {
  if err := database.DB.Where("id = ?", userID).First(&models.User{}).Error; err != nil {
    return nil, errors.New("incorrect username")
}
  claims, err := NewClaims(userID, ttl)
  if err !=nil {
    return nil, err
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    signedToken, errSign := token.SignedString([]byte(signInKey))
    if errSign != nil {
      return nil, errSign
    }

  return &Token{Value: signedToken, Claims: claims} , nil
}

func ParseToken(TokenValue string) (*tokenClaims, error) {
  token, err := jwt.ParseWithClaims(TokenValue, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("invalid signing method")
    }

    return []byte(signInKey), nil
  })
  if err != nil {
    return nil, err
  }
  
  claims, ok := token.Claims.(*tokenClaims)
  if !ok {
    return nil, errors.New("token claims are not of type")
  }

  return claims, nil
}

func UpdateTokens(userID uint, ttlAccess, ttlRefresh time.Duration) (*Tokens, error) {
  newAccessToken, errAccessToken := GenerateToken(userID, ttlAccess)
      if errAccessToken != nil {
        return nil, errAccessToken
      }
    newRefreshToken, errRefreshToken := GenerateToken(userID, ttlRefresh)
      if errRefreshToken != nil {
        return nil, errRefreshToken
      }
    return &Tokens{Access: newAccessToken, Refresh: newRefreshToken}, nil
}

func SetTokensInCookie(writer http.ResponseWriter, tokens *Tokens) {
  http.SetCookie(writer, &http.Cookie{
    Name:     "jwt_access_token",
    Value:    tokens.Access.Value,
    Path:     "/",
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    Expires:  tokens.Access.Claims.RegisteredClaims.ExpiresAt.Time,
  })

  http.SetCookie(writer, &http.Cookie{
    Name:     "jwt_refresh_token",
    Value:    tokens.Refresh.Value,
    Path:     "/",
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    Expires:  tokens.Refresh.Claims.RegisteredClaims.ExpiresAt.Time,
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