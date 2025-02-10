package utils

import (
	"context"
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func BindJSON (c *gin.Context, input interface{}) bool {
  if err:=c.ShouldBindJSON(input);err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
    return false
  }
  return true
}


func FindFavoriteMovie(userID, movieID int32) (structs.FavoriteMovie, error) {
  var favMovie structs.FavoriteMovie
  err := database.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&favMovie).Error
  return favMovie, err
}

func HardcodedUserMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    user := structs.User{
      ID:       671,
      NickName: "Vasiliy",
      UserName: "Vasya",
    }

    ctx := context.WithValue(r.Context(), "user", user)
    next.ServeHTTP(w, r.WithContext(ctx))
   })
}

func GetContextUser(ctx context.Context) (*structs.User, error) {
  user, errUser := ctx.Value("user").(structs.User)
  if !errUser {
    return nil, errors.New("User is not in context")
  }
  return &user, nil
}

func GetGenreNameByID(genreID int32, db *gorm.DB) string {
  var genre structs.Genre
  if err := db.Where("id = ?", genreID).First(&genre).Error; err != nil {
      return ""
  }
  return genre.Name
}
