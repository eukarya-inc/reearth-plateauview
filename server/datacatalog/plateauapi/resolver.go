//go:generate go run github.com/99designs/gqlgen generate --config gqlgen.yml

package plateauapi

import (
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

// Example

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	srv := plateauapi.NewSchema()

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var ErrDatacatalogUnavailable = errors.New("datacatalog is currently unavailable")

type Repo interface {
	QueryResolver
}

type Resolver struct {
	Repo Repo
}

func NewService(repo Repo) *handler.Server {
	srv := handler.NewDefaultServer(NewSchema(repo))
	return srv
}

func NewSchema(repo Repo) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{Resolvers: &Resolver{Repo: repo}})
}
