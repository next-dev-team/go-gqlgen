package cmd

import (
	"go-graph/graph/generated"
	"go-graph/graph/resolver"
	"go-graph/pkg/config"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func startServer(ctx *cli.Context) error {
	if err := config.InitDefaultServerConfig(); err != nil {
		return err
	}
	// start gorm
	if err := config.InitDefaultGormConfig(); err != nil {
		return err
	}
	conf := config.GetServerConfig()
	log.Logger = zerolog.New(&zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC822,
	}).With().Timestamp().Logger()

	logLevel, err := zerolog.ParseLevel(conf.GetLogLevel())
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLevel)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver.New(),
	}))

	http.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", conf.GetPort())
	http.ListenAndServe(":"+conf.GetPort(), nil)
	return nil
}
