package graph

import (
	"context"
	"myfavouritemovies/models"
	"myfavouritemovies/repository"
	"myfavouritemovies/service"
	"myfavouritemovies/utils"
)

func (r *mutationResolver) AddUser(ctx context.Context, nickName string, userName string, password string) (*models.User, error) {
	user := &models.User{
		NickName: nickName,
		UserName: userName,
		Password: password,
	}
	if err := repository.AddUser(user); err != nil {
		return nil, err
	}

	return user, nil
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

func (r *mutationResolver) UpdateNickName(ctx context.Context, nickName string) (*models.User, error) {
	user, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return nil, errUser
	}
	if err := repository.UpdateNickName(user, nickName); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) UpdatePassWord(ctx context.Context, password string) (*models.User, error) {
	user, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return nil, errUser
	}
	if err := repository.UpdatePassWord(user, password); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) AddFavoriteMovie(ctx context.Context, movie models.MovieInput) (*models.FavoriteMovie, error) {
	user, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return nil, errUser
	}
	favMovie, err := repository.AddFavoriteMovie(user.ID, movie)
	if err != nil {
		return nil, err
	}

	return favMovie, nil
}

func (r *mutationResolver) DeleteFavoriteMovie(ctx context.Context, favMovieID uint) (bool, error) {
	_, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return false, errUser
	}
	if err := repository.DeleteFavoriteMovie(favMovieID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ToggleWatchedStatus(ctx context.Context, favMovieID uint) (*models.FavoriteMovie, error) {
	_, errUser := utils.GetContextUser(ctx)
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
	user, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return 0, errUser
	}
	if err := repository.AddFavoriteGenre(user.ID, genreID); err != nil {
		return 0, err
	}

	return genreID, nil
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

func (r *queryResolver) GetUser(ctx context.Context) (*models.User, error) {
	user, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return nil, errUser
	}

	return user, nil
}

func (r *queryResolver) GetAllGenres(ctx context.Context) ([]*models.Genre, error) {
	genres, err := repository.GetAllGenres()
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *queryResolver) GetAllFavoriteGenres(ctx context.Context) ([]uint, error) {
	user, errUser := utils.GetContextUser(ctx)
	if errUser != nil {
		return []uint{}, errUser
	}
	favGenres, err := repository.GetFavoriteGenres(user.ID)
	if err != nil {
		return []uint{}, err
	}

	return favGenres, nil
}

func (r *queryResolver) GetFavoriteMovies(ctx context.Context) ([]*models.FavoriteMovie, error) {
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

func (r *queryResolver) GetMovieDetails(ctx context.Context, movieID uint) (*models.MovieDetails, error) {
	movieDetails, err := service.FetchMovieDetails(movieID)
	if err != nil {
		return nil, err
	}

	return movieDetails, nil
}

func (r *queryResolver) GetFilteredMovies(ctx context.Context, filter models.MovieFilter) ([]*models.Movie, error) {
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