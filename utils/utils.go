package utils

import (
	"context"
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
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
  err := database.DB.Where("id = ?", favMovieID).First(&favMovie).Error
  return favMovie, err
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
