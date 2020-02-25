package main

import (
	"log"
	"net/http"

	koverto "github.com/koverto/koverto/api"
	"github.com/koverto/koverto/internal/pkg/middleware"
	"github.com/koverto/koverto/internal/pkg/resolver"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	defaultTLSCert = "localhost+2.pem"
	defaultTLSKey  = "localhost+2-key.pem"
)

func main() {
	res, err := resolver.New()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.Handler("GET", "/", playground.Handler("Playground", "/query"))

	gqlHandler := handler.NewDefaultServer(koverto.NewExecutableSchema(*res))
	router.Handler("POST", "/query", gqlHandler)

	corsOptions := cors.Options{
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "OPTIONS", "POST"},
		AllowedOrigins: []string{"https://localhost:9000"},
	}

	chain := alice.New(
		middleware.AuthorizationHandler(res.Resolvers.(*resolver.Resolver)),
		cors.New(corsOptions).Handler,
	).Then(router)

	log.Fatal(http.ListenAndServeTLS(":8080", defaultTLSCert, defaultTLSKey, chain))
}
