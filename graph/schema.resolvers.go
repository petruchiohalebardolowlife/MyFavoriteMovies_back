package graph

import (
	"context"
	"errors"
	"myfavouritemovies/models"
	"myfavouritemovies/repository"
	"myfavouritemovies/security"
	"myfavouritemovies/service"
	"myfavouritemovies/utils"
	"net/http"
	"time"
)

func (r *mutationResolver) AddUser(ctx context.Context, nickName string, userName string, password string) (*models.User, error) {
  hash, err := security.GenerateHashPassword(password)
  if err != nil {
    return nil, err
  }
  user := &models.User{
    NickName:     nickName,
    UserName:     userName,
    PasswordHash: hash,
  }
  if err := repository.AddUser(user); err != nil {
    return nil, err
  }

  return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (bool, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.DeleteUser(userID); err != nil {
    return false, err
  }

  return true, nil
}

func (r *mutationResolver) UpdateNickName(ctx context.Context, nickName string) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, err
  }
  if err := repository.UpdateNickName(user, nickName); err != nil {
    return nil, err
  }

  return user, nil
}

func (r *mutationResolver) UpdatePassWord(ctx context.Context, password string) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, err
  }
  if err := repository.UpdatePassWord(user, password); err != nil {
    return nil, err
  }

  return user, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, signInInput models.SignInInput) (*models.User, error) {
  if err := security.SignIn(signInInput.Username, signInInput.Password); err != nil {
    return nil, err
  }
  user, errUser := repository.GetUserByUserName(signInInput.Username)
  if errUser != nil {
    return nil, errUser
  }
  accessToken, err := security.GenerateToken(user.ID, 30*time.Second)
  if err != nil {
    return nil, err
  }
  refreshToken, err := security.GenerateToken(user.ID, time.Minute)
  if err != nil {
    return nil, err
  }
  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return nil, errors.New("response writer not found")
  }
  _, errSession := repository.AddSession(&models.Session{
    ID:        refreshToken.Claims.RegisteredClaims.ID,
    UserID:    user.ID,
    ExpiresAt: refreshToken.Claims.RegisteredClaims.ExpiresAt.Time,
  })
  if errSession != nil {
    return nil, errSession
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
    return false, errUser
  }
  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return false, errors.New("response writer not found!")
  }
  request, ok := ctx.Value("httpRequest").(*http.Request)
  if !ok {
    return false, errors.New("")
  }
  refreshToken := security.TokenFromCookie(request, "jwt_refresh_token")
  claims, err := security.ParseToken(refreshToken)
  if err != nil {
    return false, err
  }

  repository.AddTokenInBlackList(claims)
  repository.DeleteSession(claims.ID)
  security.DeleteTokensFromCookie(writer)

  return true, nil
}

func (r *mutationResolver) AddFavoriteMovie(ctx context.Context, movie models.MovieInput) (*models.FavoriteMovie, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  favMovie, err := repository.AddFavoriteMovie(userID, movie)
  if err != nil {
    return nil, err
  }

  return favMovie, nil
}

func (r *mutationResolver) DeleteFavoriteMovie(ctx context.Context, favMovieID uint) (bool, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.DeleteFavoriteMovie(favMovieID); err != nil {
    return false, err
  }

  return true, nil
}

func (r *mutationResolver) ToggleWatchedStatus(ctx context.Context, favMovieID uint) (*models.FavoriteMovie, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  favMovie, err := repository.ToggleWatchedStatus(favMovieID)
  if err != nil {
    return nil, err
  }

  return favMovie, nil
}

func (r *mutationResolver) AddFavoriteGenre(ctx context.Context, genreID uint) (uint, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return 0, errUser
  }
  if err := repository.AddFavoriteGenre(userID, genreID); err != nil {
    return 0, err
  }

  return genreID, nil
}

func (r *mutationResolver) DeleteFavoriteGenre(ctx context.Context, genreID uint) (bool, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.DeleteFavoriteGenre(userID, genreID); err != nil {
    return false, err
  }

  return true, nil
}

func (r *queryResolver) GetUser(ctx context.Context) (*models.User, error) {
  UserID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  user, err := repository.GetUserByID(UserID)
  if err != nil {
    return nil, err
  }

  return user, nil
}

func (r *queryResolver) GetAllGenres(ctx context.Context) ([]*models.Genre, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  genres, err := repository.GetAllGenres()
  if err != nil {
    return nil, err
  }

  return genres, nil
}

func (r *queryResolver) GetAllFavoriteGenres(ctx context.Context) ([]uint, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return []uint{}, errUser
  }
  favGenres, err := repository.GetFavoriteGenres(userID)
  if err != nil {
    return []uint{}, err
  }

  return favGenres, nil
}

func (r *queryResolver) GetFavoriteMovies(ctx context.Context) ([]*models.FavoriteMovie, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  favMovies, err := repository.GetFavoriteMovies(userID)
  if err != nil {
    return nil, err
  }

  return favMovies, nil
}

func (r *queryResolver) GetMovieDetails(ctx context.Context, movieID uint) (*models.MovieDetails, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  movieDetails, err := service.FetchMovieDetails(movieID)
  if err != nil {
    return nil, err
  }

  return movieDetails, nil
}

func (r *queryResolver) GetFilteredMovies(ctx context.Context, filter models.MovieFilter) ([]*models.Movie, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, errUser
  }
  filteredMovies, err := service.FetchFilteredMovies(filter)
  if err != nil {
    return nil, err
  }

  return filteredMovies, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type Resolver struct{}
