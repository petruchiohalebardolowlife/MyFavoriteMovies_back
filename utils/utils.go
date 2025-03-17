package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	tokenService "myfavouritemovies/service/tokens"
	"net/http"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func GetContextUserID(ctx context.Context) (uint, error) {
  userID, ok := ctx.Value("userID").(uint)
  if !ok || userID == 0 {
    return 0, errors.New("Unauthorized")
  }

  return userID, nil
}

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    accessToken := r.Header.Get("Authorization")
    if len(accessToken) == 0 {
      next.ServeHTTP(w, r)
      return
    }
    claims, err := tokenService.Validate(accessToken)
    if err != nil {
      graphQLError := gqlerror.Errorf("Invalid token or token expired")
      graphQLError.Extensions = map[string]interface{}{
        "code": "401",
      }
      response := struct {
        Errors gqlerror.List `json:"errors"`
      }{
        Errors: gqlerror.List{graphQLError},
      }
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(http.StatusUnauthorized)

      json.NewEncoder(w).Encode(response)
      return
    }

    ctx := context.WithValue(r.Context(), "userID", claims.UserID)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func UpdateRefreshTokenInDB(currentRefreshUUID, newRefreshUUID string, newExpireAt time.Time) error {
  result := database.DB.Model(&models.Session{}).
    Where("id = ?", currentRefreshUUID).
    Updates(map[string]interface{}{
      "id":         newRefreshUUID,
      "expires_at": newExpireAt,
    })

  if result.Error != nil {
    return result.Error
  }

  if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }

  return nil
}

func HandleError(message string, code string) *gqlerror.Error {
  if code == "500" {
    log.Printf("[ERROR] %v", message)
    return &gqlerror.Error{
      Message: "Internal server error",
      Extensions: map[string]interface{}{
        "code": code,
      },
    }
  }

  return &gqlerror.Error{
    Message: message,
    Extensions: map[string]interface{}{
      "code": code,
    },
  }
}

func GetTokenFromCookie(r *http.Request) (string, error) {
  reqToken, err := r.Cookie("jwtRefresh")
  if err != nil {
    return "", err
  }

  return reqToken.Value, nil
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

func DeleteTokenFromCookie(writer http.ResponseWriter) {
  http.SetCookie(writer, &http.Cookie{
    Name:     "jwtRefresh",
    Path:     "/",
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   -1,
  })
}
