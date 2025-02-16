package utils

import (
	"context"
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"myfavouritemovies/security"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func BindJSON (c *gin.Context, input interface{}) bool {
  if err:=c.ShouldBindJSON(input);err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
    return false
  }

  return true
}


func FindFavoriteMovie(favMovieID uint) (models.FavoriteMovie, error) {
  var favMovie models.FavoriteMovie
  if err := database.DB.Where("id = ?", favMovieID).First(&favMovie).Error; err != nil{
    return models.FavoriteMovie{},err
  }
  return favMovie, nil
}

// func GetContextUser(ctx context.Context) (*models.User, error) {
//   user, errUser := ctx.Value("user").(models.User)
//   if !errUser {
//     return nil, errors.New("user is not in context")
//   }

//   return &user, nil
// }

func GetContextUserID(ctx context.Context) (uint, error) {
  userID, errUser := ctx.Value("userID").(uint)
  if !errUser {
    return 0, errors.New("user is not in context")
  }

  return userID, nil
}




func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    accessToken := security.TokenFromCookie(r,"jwt_access_token")
    if accessToken == "" {
      next.ServeHTTP(w, r)
      return
    }

    claimsAccess, err := security.ParseToken(accessToken)
    if err != nil {
      refreshToken:=security.TokenFromCookie(r, "jwt_refresh_token")
      if refreshToken == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
      security.SetTokensInCookie(w, tokens)

      ctx := context.WithValue(r.Context(), "userID", claimsRefresh.UserID)
      next.ServeHTTP(w, r.WithContext(ctx))
      return
    }

    ctx := context.WithValue(r.Context(), "userID", claimsAccess.UserID)
    next.ServeHTTP(w, r.WithContext(ctx))
   })
}

func GetUserByUserName (userName string) (*models.User, error) {
  var user *models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return nil, err
  }

  return user, nil
}
