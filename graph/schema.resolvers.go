package graph

import (
	"context"
	"myfavouritemovies/repository"
	"myfavouritemovies/service"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
)

func (r *mutationResolver) AddUser(ctx context.Context, nickName string, userName string, password string) (bool, error) {
  user := &structs.User{
    NickName: nickName,
    UserName: userName,
    Password: password,
  }
  if err := repository.AddUser(user); err != nil {
    return false, err
  }

  return true, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.DeleteUser(user.ID); err != nil {
    return false, err
  }

  return true, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, nickName *string, password *string) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.UpdateUser(user, nickName, password); err != nil {
    return false, err
  }

  return true, nil
}

func (r *mutationResolver) AddFavoriteMovie(ctx context.Context, movie structs.MovieInput) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.AddFavoriteMovie(user.ID, movie); err != nil {
    return false, err
  }
  return true, nil
}

func (r *mutationResolver) DeleteFavoriteMovie(ctx context.Context, movieID uint) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.DeleteFavoriteMovie(user.ID, movieID); err != nil {
    return false, err
  }
  return true, nil
}

func (r *mutationResolver) ToggleWatchedStatus(ctx context.Context, movieID uint) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.ToggleWatchedStatus(user.ID, movieID); err != nil {
    return false, err
  }
  return true, nil
}

func (r *mutationResolver) AddFavoriteGenre(ctx context.Context, genreID uint) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.AddFavoriteGenre(user.ID, genreID); err != nil {
    return false, err
  }
  return true, nil
}

func (r *mutationResolver) DeleteFavoriteGenre(ctx context.Context, genreID uint) (bool, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return false, errUser
  }
  if err := repository.DeleteFavoriteGenre(user.ID, genreID); err != nil {
    return false, err
  }
  return true, nil
}

func (r *queryResolver) GetUser(ctx context.Context) (*structs.User, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return nil, errUser
  }
  return user, nil
}

func (r *queryResolver) GetAllGenres(ctx context.Context) ([]*structs.Genre, error) {
  genres, err := repository.GetAllGenres()
  if err != nil {
    return nil, err
  }

  return genres, nil
}

func (r *queryResolver) GetAllFavoriteGenres(ctx context.Context) ([]*structs.FavoriteGenre, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return nil, errUser
  }
  favGenres, err := repository.GetFavoriteGenres(user.ID)
  if err != nil {
    return nil, err
  }
  return favGenres, nil
}

func (r *queryResolver) GetFavoriteMovies(ctx context.Context) ([]*structs.FavoriteMovie, error) {
  user, errUser := utils.GetContextUser(ctx)
  if errUser != nil {
    return nil, errUser
  }
  favMovies, err := repository.GetFavoriteMovies(user.ID)
  if err != nil {
    return nil, err
  }
  return favMovies, nil
}

func (r *queryResolver) GetMovieDetails(ctx context.Context, movieID uint) (*structs.MovieDetails, error) {
  movieDetails, err := service.FetchMovieDetails(movieID)
  if err != nil {
    return nil, err
  }

  return movieDetails, nil
}

func (r *queryResolver) GetFilteredMovies(ctx context.Context, filter structs.MovieFilter) ([]*structs.Movie, error) {
  filteredMovies, err := service.FetchFilteredMovies(filter)
  if err != nil {
    return nil, err
  }
  return filteredMovies, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
type Resolver struct{}
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
