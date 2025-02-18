package utils

import (
	"context"
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"myfavouritemovies/repository"
	"myfavouritemovies/security"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func GetContextUserID(ctx context.Context) (uint, error) {
  userID, errUser := ctx.Value("userID").(uint)
  if !errUser {
    return 0, errors.New("Unathorized")
  }

  return userID, nil
}

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    accessToken:= security.TokenFromCookie(r,"jwt_access_token")

    if accessToken == "" {
      CheckRefreshToken(w, r, next)
      return
    }

    claimsAccess, err := security.ParseToken(accessToken)
    if err != nil {
      CheckRefreshToken(w, r, next)
      return
    }

    ctx := context.WithValue(r.Context(), "userID", claimsAccess.UserID)
    next.ServeHTTP(w, r.WithContext(ctx))
   })
}

func CheckRefreshToken(w http.ResponseWriter, r *http.Request, next http.Handler) {
  refreshToken:=security.TokenFromCookie(r, "jwt_refresh_token")
  if refreshToken == "" {
    next.ServeHTTP(w,r)
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
    http.Error(w, "Failed to add refresh token in blacklist",http.StatusUnauthorized)
    return
  }

  if err := UpdateRefreshTokenInDB(claimsRefresh.RegisteredClaims.ID, tokens.Refresh.Claims.ID, time.Now(), tokens.Refresh.Claims.ExpiresAt.Time); err != nil {
    http.Error(w, "Session Not Found", http.StatusUnauthorized)
    return
  }

  security.SetTokensInCookie(w, tokens)
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
