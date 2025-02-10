package main

import (
	"fmt"
	"log"
	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
	"myfavouritemovies/graph"
	"myfavouritemovies/repository"
	"myfavouritemovies/service"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)


func main() {
  db := database.InitDB()
  if db == nil {
    log.Fatal("Failed to initialize the database.")
  }
  fmt.Println("Database tables created successfully!")
  resolver := &graph.Resolver{}

  genres, err := service.FetchGenres()
  if err != nil {
    log.Fatalf("Failed to fetch genres: %v", err)
  }
  if err := repository.SaveGenresToDB(genres); err != nil {
    log.Fatalf("Failed to add genres: %v", err)
  }

  srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

  srv.AddTransport(transport.Options{})
  srv.AddTransport(transport.GET{})
  srv.AddTransport(transport.POST{})

  srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

  srv.Use(extension.Introspection{})
  srv.Use(extension.AutomaticPersistedQuery{
    Cache: lru.New[string](100),
  })

  http.Handle("/", utils.HardcodedUserMiddleware(playground.Handler("GraphQL playground", "/query")))
  http.Handle("/query", utils.HardcodedUserMiddleware(srv))

  log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.SRVR_PORT)
  log.Fatal(http.ListenAndServe(":"+config.SRVR_PORT, nil))
}
