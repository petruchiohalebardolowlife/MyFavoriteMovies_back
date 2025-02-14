package utils

import (
	"context"
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"myfavouritemovies/security"
	"net/http"

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

func HardcodedUserMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    user := models.User{
      ID:       671,
      NickName: "Vasiliy",
      UserName: "Vasya",
    }

    ctx := context.WithValue(r.Context(), "user", user)
    next.ServeHTTP(w, r.WithContext(ctx))
   })
}

func GetContextUser(ctx context.Context) (*models.User, error) {
  user, errUser := ctx.Value("user").(models.User)
  if !errUser {
    return nil, errors.New("user is not in context")
  }

  return &user, nil
}


func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    token := security.TokenFromHTTPRequest(r)
    ctx := context.WithValue(r.Context(), "httpResponseWriter", w)
    if token == "" {
      next.ServeHTTP(w, r.WithContext(ctx))
      return
    }

    userID, err := security.ParseToken(token)
    if err != nil {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }
    var user models.User
    if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }
    ctx = context.WithValue(r.Context(), "user", user)
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
