package graph

import (
	"context"
	"log"
	"myfavouritemovies/repository"
	"myfavouritemovies/structs"
)

type Resolver struct{}

func (r *Resolver) getUser(ctx context.Context, id int32) (*structs.User, error) {
	return &structs.User{ID: id, Username: "John Doe", FavoriteMovies: []*structs.FavoriteMovie{}}, nil
}

func(r *Resolver) getAllGenres(ctx context.Context) ([]*structs.Genre, error) {
  genres, err := repository.GetAllGenres()
  if err != nil {
      log.Println("Ошибка запроса в БД:", err)
      return nil, err
  }
  genrePointers := make([]*structs.Genre, len(genres))
    for i := range genres {
        genrePointers[i] = &genres[i]
    }
  return genrePointers, nil
}

func (r *Resolver) Query_movies(ctx context.Context, filter *structs.MovieFilter) ([]*structs.Movie, error) {
	return []*structs.Movie{
		{ID: 1, Title: "Inception", VoteAverage: 8.8},
	}, nil
}

func (r *Resolver) Mutation_addFavoriteMovie(ctx context.Context, userID string, movieID string) (*structs.User, error) {
	return &structs.User{ID: 168, Username: "John Doe"}, nil
}