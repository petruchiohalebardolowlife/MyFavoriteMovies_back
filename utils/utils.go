package utils

import (
	"context"
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"myfavouritemovies/security"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func FindFavoriteMovie(favMovieID uint) (models.FavoriteMovie, error) {
  var favMovie models.FavoriteMovie
  if err := database.DB.Where("id = ?", favMovieID).First(&favMovie).Error; err != nil{
    return models.FavoriteMovie{},err
  }
  
  return favMovie, nil
}

func GetContextUserID(ctx context.Context) (uint, error) {
  userID, errUser := ctx.Value("userID").(uint)
  if !errUser {
    return 0, errors.New("user is not in context")
  }

  return userID, nil
}

func GetUserByUserName (userName string) (*models.User, error) {
  var user *models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return nil, err
  }

  return user, nil
}

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    accessToken := security.TokenFromCookie(r,"jwt_access_token")

    fingerprint, err := security.FingerPrintFromHTTPRequest(r)
    if err != nil {
      return
    }
    ctx := context.WithValue(r.Context(), "Fingerprint", fingerprint)
    
    if accessToken == "" {
      CheckRefreshToken(w, r.WithContext(ctx), next)
      return
    }

    claimsAccess, err := security.ParseToken(accessToken)
    if err != nil {
      CheckRefreshToken(w, r.WithContext(ctx), next)
      return
    }

    ctx = context.WithValue(r.Context(), "userID", claimsAccess.UserID)
    ctx = context.WithValue(ctx, "Fingerprint", fingerprint)
    next.ServeHTTP(w, r.WithContext(ctx))
   })
}

func CheckRefreshToken(w http.ResponseWriter, r *http.Request, next http.Handler) {
  refreshToken:=security.TokenFromCookie(r, "jwt_refresh_token")
  if refreshToken == "" {
    next.ServeHTTP(w, r)
    return 
  }

  claimsRefresh, err := security.ParseToken(refreshToken)
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return 
  }
  

  tokens, errTokens := security.UpdateTokens(claimsRefresh.UserID, 30*time.Second, time.Minute)
  if errTokens != nil {
    http.Error(w, "Failed to generates tokens", http.StatusUnauthorized)
    return
  }

  fingerprint, ok := r.Context().Value("Fingerprint").(string)
  if !ok  {
    http.Error(w, "Fingerprint not found", http.StatusUnauthorized)
    return
  }

  if err := UpdateRefreshTokenInDB(claimsRefresh.RegisteredClaims.ID, tokens.Refresh.Claims.ID, fingerprint, time.Now(), tokens.Refresh.Claims.ExpiresAt.Time); err != nil {
    http.Error(w, "Session Not Found", http.StatusUnauthorized)
    return
  }

  security.SetTokensInCookie(w, tokens)

  ctx := context.WithValue(r.Context(), "userID", claimsRefresh.UserID)
  next.ServeHTTP(w, r.WithContext(ctx))
}

func UpdateRefreshTokenInDB(refreshUUID, newRefreshUUID, fingerprint string, now, newExpireAt time.Time) error {
  result := database.DB.Model(&models.Session{}).
    Where("id = ? AND fingerprint = ? AND expires_at > ?", refreshUUID, fingerprint, now).
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
