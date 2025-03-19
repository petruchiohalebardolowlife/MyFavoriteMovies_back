// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	config "myfavouritemovies/configs"
// 	"myfavouritemovies/database"
// 	"myfavouritemovies/graph"
// 	"myfavouritemovies/repository"
// 	service "myfavouritemovies/service/tmdb"
// 	"myfavouritemovies/utils"
// 	"net/http"
// 	"time"

// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/handler/extension"
// 	"github.com/99designs/gqlgen/graphql/handler/lru"
// 	"github.com/99designs/gqlgen/graphql/handler/transport"
// 	"github.com/99designs/gqlgen/graphql/playground"
// 	"github.com/rs/cors"
// 	"github.com/vektah/gqlparser/v2/ast"
// )

// func main() {
//   db := database.InitDB()
//   if db != nil {
//     fmt.Println("Database tables created successfully!")
//   } else {
//     log.Fatal("Failed to initialize the database.")
//   }
//   repository.CleanExpiredTokens(24 * time.Hour)

//   resolver := &graph.Resolver{}
//   genres, err := service.FetchGenres()
//   if err != nil {
//     log.Fatalf("Failed to fetch genres: %v", err)
//   }
//   if err := repository.SaveGenresToDB(genres); err != nil {
//     log.Fatalf("Failed to add genres: %v", err)
//   }
//   srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

//   srv.AddTransport(transport.Options{})
//   srv.AddTransport(transport.GET{})
//   srv.AddTransport(transport.POST{})

//   srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

//   srv.Use(extension.Introspection{})
//   srv.Use(extension.AutomaticPersistedQuery{
//     Cache: lru.New[string](100),
//   })

//   if config.APP_ENV == "development" {
//     http.Handle("/", playground.Handler("GraphQL playground", "/query"))
//   } else {
//     http.Handle("/", http.NotFoundHandler())
//   }
//   c := cors.New(cors.Options{
//     AllowedOrigins:   []string{"http://localhost:5173"},
//     AllowCredentials: true,
//     AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
//     AllowedHeaders:   []string{"Authorization", "Content-Type"},
// })
//   http.Handle("/query", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     ctx := context.WithValue(r.Context(), "httpResponseWriter", w)
//     ctx = context.WithValue(ctx, "httpRequest", r)
//     utils.Middleware(srv).ServeHTTP(w, r.WithContext(ctx))
// })))

//   log.Fatal(http.ListenAndServe(":"+config.SRVR_PORT, nil))
// }

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
	"myfavouritemovies/graph"
	"myfavouritemovies/repository"
	service "myfavouritemovies/service/tmdb"
	"myfavouritemovies/utils"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       string
}

func main() {
	// Инициализация базы данных
	db := database.InitDB()
	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}

	// Чистим устаревшие токены
	repository.CleanExpiredTokens(24 * time.Hour)

	// Настройка GraphQL resolver
	resolver := &graph.Resolver{}
	genres, err := service.FetchGenres()
	if err != nil {
		log.Fatalf("Failed to fetch genres: %v", err)
	}
	if err := repository.SaveGenresToDB(genres); err != nil {
		log.Fatalf("Failed to add genres: %v", err)
	}

	// Создание GraphQL сервера
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Настройка транспортеров
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

  srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

  srv.Use(extension.Introspection{})
  srv.Use(extension.AutomaticPersistedQuery{
    Cache: lru.New[string](100),
  })

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	// Плейграунд GraphQL
	if config.APP_ENV == "development" {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	} else {
		http.Handle("/", http.NotFoundHandler())
	}

	// Обработка запросов /query
	http.Handle("/query", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запроса
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

		// Логирование заголовков
		for name, values := range r.Header {
			fmt.Printf("Header: %s = %v\n", name, values)
		}

		// Логирование тела запроса, если POST
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Error reading body:", err)
			} else {
				fmt.Println("Request body:", string(body)) // Логируем тело запроса
				// Сохраняем тело обратно в r.Body для дальнейшей обработки
				r.Body = io.NopCloser(bytes.NewReader(body))
			}
		}

		// Логируем, что сервер приступает к обработке запроса
		fmt.Println("Processing request...")

		// Прокси для логирования ответа
		responseRecorder := &responseRecorder{ResponseWriter: w}

		// Применяем свой middleware перед вызовом GraphQL обработчика
		// Предполагаем, что `utils.Middleware(srv)` - это middleware, которое ты хочешь использовать
		fmt.Println("Handling request directly...")

		// Применяем твой middleware к запросу
		// Здесь мы предполагаем, что твой middleware правильно обрабатывает запросы
		utils.Middleware(srv).ServeHTTP(responseRecorder, r)
		// Для отладки заменим его на прямой вызов, чтобы проверить логи:
		// srv.ServeHTTP(responseRecorder, r)

		// Логирование ответа после обработки запроса
		fmt.Printf("Response Status Code: %d\n", responseRecorder.statusCode)
		fmt.Printf("Response Body: %s\n", responseRecorder.body)

		// Проверка на пустое тело ответа
		if responseRecorder.statusCode == 0 {
			fmt.Println("Warning: Response body is empty or was not set.")
		}

		// Если ошибка при обработке запроса, логируем это
		if responseRecorder.statusCode >= 400 {
			fmt.Printf("Error: Request returned an error with status code %d\n", responseRecorder.statusCode)
		}

		// Завершаем логирование
		fmt.Println("Request handling completed.")
	})))

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":"+config.SRVR_PORT, nil))
}
