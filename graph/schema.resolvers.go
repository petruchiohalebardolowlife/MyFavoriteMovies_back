package graph

import (
	"context"
	"myfavouritemovies/models"
	"myfavouritemovies/repository"
	"myfavouritemovies/security"
	tmdbService "myfavouritemovies/service/tmdb"
	tokenService "myfavouritemovies/service/tokens"
	"myfavouritemovies/utils"
	"net/http"
)

func (r *mutationResolver) AddUser(ctx context.Context, nickName string, userName string, password string) (*models.User, error) {
  hash, err := security.GenerateHashPassword(password)
  if err != nil {
    return nil, utils.HandleError("Failed to generate hash password", "500")
  }
  user := &models.User{
    NickName:     nickName,
    UserName:     userName,
    PasswordHash: hash,
  }
  if err := repository.AddUser(user); err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (bool, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, utils.HandleError("Unauthorized", "401")
  }
  if err := repository.DeleteUser(userID); err != nil {
    return false, utils.HandleError("DB Error", "500")
  }

  return true, nil
}

func (r *mutationResolver) UpdateNickName(ctx context.Context, nickName string) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }
  if err := repository.UpdateNickName(user, nickName); err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return user, nil
}

func (r *mutationResolver) UpdatePassWord(ctx context.Context, password string) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }
  if err := repository.UpdatePassWord(user, password); err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return user, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, signInInput models.SignInInput) (*models.SignInResponse, error) {
  if err := security.SignIn(signInInput.Username, signInInput.Password); err != nil {
    return nil, utils.HandleError("Invalid username or password", "401")
  }
  user, errUser := repository.GetUserByUserName(signInInput.Username)
  if errUser != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  accessToken, refreshToken, err := tokenService.Generate(user.ID)
  if err != nil {
    return nil, utils.HandleError("Failed to generate tokens", "500")
  }

  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return nil, utils.HandleError("httpResponseWriter not found", "500")
  }

  refreshTokenClaims, err := tokenService.Validate(refreshToken)
  if err != nil {
    return nil, utils.HandleError(err.Error(), "401")
  }
  _, errSession := repository.AddSession(&models.Session{
    ID:        refreshTokenClaims.ID,
    UserID:    user.ID,
    ExpiresAt: refreshTokenClaims.ExpiresAt.Time,
  })
  if errSession != nil {
    return nil, utils.HandleError("DB Error", "500")
  }
  utils.SetTokenInCookie(writer, refreshToken)

  return &models.SignInResponse{User: user, Token: accessToken}, nil
}

func (r *mutationResolver) LogOut(ctx context.Context) (bool, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, utils.HandleError("Unauthorized", "401")
  }
  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return false, utils.HandleError("httpResponseWriter not found", "500")
  }
  request, ok := ctx.Value("httpRequest").(*http.Request)
  if !ok {
    return false, utils.HandleError("httpRequest not found", "500")
  }
  refreshToken, errRefreshToken := utils.GetTokenFromCookie(request)
  if errRefreshToken != nil {
    return false, utils.HandleError("Unauthorized", "401")
  }
  claims, err := tokenService.Validate(refreshToken)
  if err != nil {
    return false, utils.HandleError(err.Error(), "401")
  }

  if err := repository.AddTokenInBlackList(claims); err != nil {
    return false, utils.HandleError("DB Error", "500")
  }

  if err := repository.DeleteSession(claims.ID); err != nil {
    return false, utils.HandleError("DB Error", "500")
  }

  utils.DeleteTokenFromCookie(writer)

  return true, nil
}

func (r *mutationResolver) AddFavoriteMovie(ctx context.Context, movie models.MovieInput) (*models.FavoriteMovie, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  favMovie, err := repository.AddFavoriteMovie(userID, movie)
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return favMovie, nil
}

func (r *mutationResolver) DeleteFavoriteMovie(ctx context.Context, favMovieID uint) (bool, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, utils.HandleError("Unauthorized", "401")
  }
  if err := repository.DeleteFavoriteMovie(favMovieID); err != nil {
    return false, utils.HandleError("DB Error", "500")
  }

  return true, nil
}

func (r *mutationResolver) ToggleWatchedStatus(ctx context.Context, favMovieID uint) (*models.FavoriteMovie, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  favMovie, err := repository.ToggleWatchedStatus(favMovieID)
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return favMovie, nil
}

func (r *mutationResolver) AddFavoriteGenre(ctx context.Context, genreID uint) (uint, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return 0, utils.HandleError("Unauthorized", "401")
  }
  if err := repository.AddFavoriteGenre(userID, genreID); err != nil {
    return 0, utils.HandleError("DB Error", "500")
  }

  return genreID, nil
}

func (r *mutationResolver) DeleteFavoriteGenre(ctx context.Context, genreID uint) (bool, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return false, utils.HandleError("Unauthorized", "401")
  }
  if err := repository.DeleteFavoriteGenre(userID, genreID); err != nil {
    return false, utils.HandleError("DB Error", "500")
  }

  return true, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (string, error) {
  request, ok := ctx.Value("httpRequest").(*http.Request)
  if !ok {
    return "", utils.HandleError("httpRequest not found", "500")
  }
  writer, ok := ctx.Value("httpResponseWriter").(http.ResponseWriter)
  if !ok {
    return "", utils.HandleError("httpResponseWriter not found", "500")
  }
  currentRefreshTokenCookie, err := request.Cookie("jwtRefresh")
  if err != nil {
    return "", utils.HandleError("No refresh token in cookies", "401")
  }
  currentClaimsRefresh, err := tokenService.Validate(currentRefreshTokenCookie.Value)
  if err != nil {
    return "", utils.HandleError(err.Error(), "401")
  }
  newAccessToken, newRefreshToken, err := tokenService.Refresh(currentRefreshTokenCookie.Value)
  if err != nil {
    return "", utils.HandleError(err.Error(), "401")
  }
  newClaimsRefresh, err := tokenService.Validate(newRefreshToken)
  if err != nil {
    return "", utils.HandleError(err.Error(), "401")
  }
  if err := repository.AddTokenInBlackList(currentClaimsRefresh); err != nil {
    return "", utils.HandleError("DB Error", "500")
  }
  if err := utils.UpdateRefreshTokenInDB(currentClaimsRefresh.ID, newClaimsRefresh.ID, newClaimsRefresh.ExpiresAt.Time); err != nil {
    return "", utils.HandleError("DB Error", "500")
  }
  utils.SetTokenInCookie(writer, newRefreshToken)

  return newAccessToken, nil
}

func (r *queryResolver) GetUser(ctx context.Context) (*models.User, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  user, err := repository.GetUserByID(userID)
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return user, nil
}

func (r *queryResolver) GetAllGenres(ctx context.Context) ([]*models.Genre, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  genres, err := repository.GetAllGenres()
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return genres, nil
}

func (r *queryResolver) GetAllFavoriteGenres(ctx context.Context) ([]uint, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return []uint{}, utils.HandleError("Unauthorized", "401")
  }
  favGenres, err := repository.GetFavoriteGenres(userID)
  if err != nil {
    return []uint{}, utils.HandleError("DB Error", "500")
  }

  return favGenres, nil
}

func (r *queryResolver) GetFavoriteMovies(ctx context.Context) ([]*models.FavoriteMovie, error) {
  userID, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  favMovies, err := repository.GetFavoriteMovies(userID)
  if err != nil {
    return nil, utils.HandleError("DB Error", "500")
  }

  return favMovies, nil
}

func (r *queryResolver) GetMovieDetails(ctx context.Context, movieID uint) (*models.MovieDetails, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  movieDetails, err := tmdbService.FetchMovieDetails(movieID)
  if err != nil {
    return nil, utils.HandleError("Service Unavailable", "503")
  }

  return movieDetails, nil
}

func (r *queryResolver) GetFilteredMovies(ctx context.Context, filter models.MovieFilter) ([]*models.Movie, error) {
  _, errUser := utils.GetContextUserID(ctx)
  if errUser != nil {
    return nil, utils.HandleError("Unauthorized", "401")
  }
  filteredMovies, err := tmdbService.FetchFilteredMovies(filter)
  if err != nil {
    return nil, utils.HandleError("Service Unavailable", "503")
  }

  return filteredMovies, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type Resolver struct{}
