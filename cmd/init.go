package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type AppInfo struct {
	Commit     string
	Build      int
	Name       string
	Version    string
	Usage      string
	CommitHash string
	CompiledAt time.Time
}

func Init(appInfo *AppInfo) {
	app := &cli.App{
		Name:    appInfo.Name,
		Version: fmt.Sprintf("%v-%d-%s", appInfo.Version, appInfo.Build, appInfo.Commit),
		Usage:   appInfo.Usage,
		Action:  startServer,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.toml",
				Usage: "load configuration from `FILE` in sub directory configs",
			},
			&cli.BoolFlag{
				Name:  "cpuprofile",
				Usage: "enable cpu profile",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(fmt.Errorf("execute failed: %v", err))
	}
}
