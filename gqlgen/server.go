package gqlgen

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gqlgen/graph"
	"net/http"
)

func GraphQLPlaygroundHandler(route string) (*handler.Server, http.HandlerFunc) {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	return srv, playground.Handler("GraphQL playground", route)
}
