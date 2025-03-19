// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type FilteredMoviesResponse struct {
	Page    uint     `json:"page"`
	Results []*Movie `json:"results,omitempty"`
}

type GetFavoriteMoviesResponse struct {
	Page       uint             `json:"page"`
	Results    []*FavoriteMovie `json:"results,omitempty"`
	TotalPages uint             `json:"totalPages"`
}

type MovieFilter struct {
	GenreIDs   []uint   `json:"genreIDs,omitempty"`
	Popularity *float64 `json:"popularity,omitempty"`
	Year       *uint    `json:"year,omitempty"`
	Page       *uint    `json:"page,omitempty"`
}

type MovieInput struct {
	MovieID     uint    `json:"movieID"`
	Title       string  `json:"title"`
	PosterPath  string  `json:"posterPath"`
	VoteAverage float64 `json:"voteAverage"`
	GenreIDs    []uint  `json:"genreIDs"`
}

type Mutation struct {
}

type Query struct {
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
