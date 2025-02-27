package utils

import (
	"context"
	"errors"
	"log"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"myfavouritemovies/repository"
	"myfavouritemovies/security"
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

    userID, err := tokenService.Validate(accessToken)
    if err != nil {
      http.Error(w, "Token expired", http.StatusUnauthorized)
    }

    ctx := context.WithValue(r.Context(), "userID", userID)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func CheckRefreshToken(w http.ResponseWriter, r *http.Request, next http.Handler) {
  refreshToken, err := security.TokenFromCookie(r)
  if err != nil {
    return
  }
  if refreshToken == "" {
    next.ServeHTTP(w, r)
    return
  }

  claimsRefresh, err := security.ParseToken(refreshToken)
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  if err := repository.CheckTokenInBlackList(claimsRefresh.ID); err != nil {
    http.Error(w, "Your refresh token in blacklist", http.StatusUnauthorized)
    return
  }

  tokens, errTokens := security.UpdateTokens(claimsRefresh.UserID, 15*time.Minute, 60*24*time.Hour)
  if errTokens != nil {
    http.Error(w, "Failed to generates tokens", http.StatusUnauthorized)
    return
  }

  if err := repository.AddTokenInBlackList(claimsRefresh); err != nil {
    http.Error(w, "Failed to add refresh token in blacklist", http.StatusUnauthorized)
    return
  }

  if err := UpdateRefreshTokenInDB(claimsRefresh.RegisteredClaims.ID, tokens.Refresh.Claims.ID, time.Now(), tokens.Refresh.Claims.ExpiresAt.Time); err != nil {
    http.Error(w, "Session Not Found", http.StatusUnauthorized)
    return
  }

  // security.SetTokensInCookie(w, tokens)
  ctx := context.WithValue(r.Context(), "userID", claimsRefresh.UserID)
  next.ServeHTTP(w, r.WithContext(ctx))
}

func UpdateRefreshTokenInDB(refreshUUID, newRefreshUUID string, now, newExpireAt time.Time) error {
  result := database.DB.Model(&models.Session{}).
    Where("id = ? AND expires_at > ?", refreshUUID, now).
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
