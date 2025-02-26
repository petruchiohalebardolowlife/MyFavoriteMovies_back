package graph

import (
	"context"
	"log"
	"myfavouritemovies/models"
	"myfavouritemovies/repository"
	"myfavouritemovies/security"
	"myfavouritemovies/service"
	"myfavouritemovies/utils"
	"net/http"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func handleError(message string, code string) *gqlerror.Error {
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

func (r *mutationResolver) AddUser(ctx context.Context, nickName string, userName string, password string) (*models.User, error) {
  hash, err := security.GenerateHashPassword(password)
  if err != nil {
    return nil, handleError("Failed to generate hash password", "500")
  }
  user := &models.User{
    NickName: nickName,
    UserName: userName,
    PasswordHash: hash,
  }
  if err := repository.AddUser(user); err != nil {
    return nil, handleError("DB Error", "500")
  }

  return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (bool, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, handleError("Unauthorized", "401")
  }
  if err := repository.DeleteUser(userID); err != nil {
    return false, handleError("DB Error", "500")
  }

  return true, nil
}

func (r *mutationResolver) UpdateNickName(ctx context.Context, nickName string) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, handleError("DB Error", "500")
  }
  if err := repository.UpdateNickName(user, nickName); err != nil {
    return nil, handleError("DB Error", "500")
  }

  return user, nil
}

func (r *mutationResolver) UpdatePassWord(ctx context.Context, password string) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, handleError("DB Error", "500")
  }
  if err := repository.UpdatePassWord(user, password); err != nil {
    return nil, handleError("DB Error", "500")
  }

  return user, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, signInInput models.SignInInput) (*models.User, error) {
  if err := security.SignIn(signInInput.Username, signInInput.Password); err != nil {
    return nil, handleError("Invalid username or password", "401")
  }
  user, errUser := repository.GetUserByUserName(signInInput.Username)
  if errUser != nil {
    return nil, handleError("Failed to get user from DB", "500")
  }
  accessToken, err := security.GenerateToken(user.ID, 15*time.Minute)
  if err != nil {
    return nil, handleError("Failed to generate token", "500")
  }
  refreshToken, err := security.GenerateToken(user.ID, 60*24*time.Hour)
  if err != nil {
    return nil, handleError("Failed to generate token", "500")
  }
  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return nil, handleError("httpResponseWriter not found","500")
  }
  _, errSession := repository.AddSession(&models.Session{
    ID:        refreshToken.Claims.RegisteredClaims.ID,
    UserID:    user.ID,
    ExpiresAt: refreshToken.Claims.RegisteredClaims.ExpiresAt.Time,
  })
  if errSession != nil {
    return nil, handleError("DB Error","500")
  }
  security.SetTokensInCookie(writer, &models.Tokens{
    Access:  accessToken,
    Refresh: refreshToken,
  })

  return user, nil
}

func (r *mutationResolver) LogOut(ctx context.Context) (bool, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, handleError("Unauthorized", "401")
  }
  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return false, handleError("httpResponseWriter not found", "500")
  }
  request, ok := ctx.Value("httpRequest").(*http.Request)
  if !ok {
    return false, handleError("httpRequest not found", "500")
  }
  refreshToken, errRefreshToken := security.TokenFromCookie(request)
  if errRefreshToken != nil {
    return false, handleError("Unauthorized", "401")
  }
  claims, err := security.ParseToken(refreshToken)
  if err != nil {
    return false, handleError("Refresh token expired", "401")
  }

  repository.AddTokenInBlackList(claims)
  repository.DeleteSession(claims.ID)
  security.DeleteTokensFromCookie(writer)

  return true, nil
}

func (r *mutationResolver) AddFavoriteMovie(ctx context.Context, movie models.MovieInput) (*models.FavoriteMovie, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  favMovie, err := repository.AddFavoriteMovie(userID, movie)
  if err != nil {
    return nil, handleError("DB Error","500")
  }

  return favMovie, nil
}

func (r *mutationResolver) DeleteFavoriteMovie(ctx context.Context, favMovieID uint) (bool, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, handleError("Unauthorized", "401")
  }
  if err := repository.DeleteFavoriteMovie(favMovieID); err != nil {
    return false, handleError("DB Error","500")
  }

  return true, nil
}

func (r *mutationResolver) ToggleWatchedStatus(ctx context.Context, favMovieID uint) (*models.FavoriteMovie, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  favMovie, err := repository.ToggleWatchedStatus(favMovieID)
  if err != nil {
    return nil, handleError("DB Error","500")
  }

  return favMovie, nil
}

func (r *mutationResolver) AddFavoriteGenre(ctx context.Context, genreID uint) (uint, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return 0, handleError("Unauthorized", "401")
  }
  if err := repository.AddFavoriteGenre(userID, genreID); err != nil {
    return 0, handleError("DB Error","500")
  }

  return genreID, nil
}

func (r *mutationResolver) DeleteFavoriteGenre(ctx context.Context, genreID uint) (bool, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, handleError("Unauthorized", "401")
  }
  if err := repository.DeleteFavoriteGenre(userID, genreID); err != nil {
    return false, handleError("DB Error","500")
  }

  return true, nil
}

func (r *queryResolver) GetUser(ctx context.Context) (*models.User, error) {
  UserID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  user, err := repository.GetUserByID(UserID)
  if err != nil {
    return nil, handleError("DB Error","500")
  }

  return user, nil
}

func (r *queryResolver) GetAllGenres(ctx context.Context) ([]*models.Genre, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  genres, err := repository.GetAllGenres()
  if err != nil {
    return nil, handleError("DB Error","500")
  }

  return genres, nil
}

func (r *queryResolver) GetAllFavoriteGenres(ctx context.Context) ([]uint, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return []uint{}, handleError("Unauthorized", "401")
  }
  favGenres, err := repository.GetFavoriteGenres(userID)
  if err != nil {
    return []uint{}, handleError("DB Error","500")
  }

  return favGenres, nil
}

func (r *queryResolver) GetFavoriteMovies(ctx context.Context) ([]*models.FavoriteMovie, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  favMovies, err := repository.GetFavoriteMovies(userID)
  if err != nil {
    return nil, handleError("DB Error","500")
  }

  return favMovies, nil
}

func (r *queryResolver) GetMovieDetails(ctx context.Context, movieID uint) (*models.MovieDetails, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  movieDetails, err := service.FetchMovieDetails(movieID)
  if err != nil {
    return nil, handleError("Service Unavailable", "503")
  }

  return movieDetails, nil
}

func (r *queryResolver) GetFilteredMovies(ctx context.Context, filter models.MovieFilter) ([]*models.Movie, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, handleError("Unauthorized", "401")
  }
  filteredMovies, err := service.FetchFilteredMovies(filter)
  if err != nil {
    return nil, handleError("Service Unavailable", "503")
  }

  return filteredMovies, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type Resolver struct{}
