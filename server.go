package main

import (
	"example/graph/resolvers"
	"example/graph/generated"
	"log"
	"net/http"
	"os"
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/99designs/gqlgen/graphql"
	"example/graph/logengine"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	postgres "example/graph/postgres"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load(".dev.env")

	if err != nil {
		log.Println("error loading env", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	postgres.InitDbPool()

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "*"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		goc := graphql.GetOperationContext(ctx)
		if goc.OperationName != "IntrospectionQuery" {
			if goc.Operation.Operation == "query" {
				logengine.GetTelemetryClient().TrackEvent(string(goc.Operation.Operation) + " " + goc.RawQuery)
			} else {
				logengine.GetTelemetryClient().TrackEvent(goc.RawQuery)
			}
		}
		return next(ctx)
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
