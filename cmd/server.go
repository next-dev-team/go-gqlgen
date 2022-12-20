package cmd

import (
	"go-graph/graph/generated"
	"go-graph/graph/resolver"
	"go-graph/pkg/config"
	"go-graph/pkg/splitlog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	zrm *splitlog.ZeroLogRotateManager
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
	logLevel, err := zerolog.ParseLevel(conf.GetLogLevel())
	if err != nil {
		return err
	}
	initSplitLog(logLevel)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver.New(),
	}))

	http.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", conf.GetPort())
	http.ListenAndServe(":"+conf.GetPort(), nil)
	return nil
}

func initSplitLog(logLevel zerolog.Level) {
	yy, mm, dd := time.Now().Date()
	tomorrowMidNight := time.Date(yy, mm, dd+1, 0, 0, 0, 0, time.Local)
	zerolog.SetGlobalLevel(logLevel)
	zrm = &splitlog.ZeroLogRotateManager{
		SplitLogLevel:   true,
		DefaultLogLevel: zerolog.InfoLevel,
		Dir:             "logs",
		FormatFilename:  "2006_01_02",
		FirstRotation:   time.Until(tomorrowMidNight),
		RotateDuration:  24 * time.Hour,
		Console:         true,
	}
	zrm.Init()
	log.Logger = zrm.Logger
	log.Info().Msg("log rotation has started")
}
