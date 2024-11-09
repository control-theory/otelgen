package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var logger *zap.Logger

func initLogger(c *cli.Context) error {

	var cfg zap.Config
	var err error

	switch c.String("log-level") {
	case "debug":
		cfg = zap.NewDevelopmentConfig()
	case "info":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level.SetLevel(zap.InfoLevel)
	case "warn":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level.SetLevel(zap.WarnLevel)
	case "error":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level.SetLevel(zap.ErrorLevel)
	default:
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(zap.WarnLevel)
	}
	logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync() // nolint: errcheck

	return err
}

func New(version, commit, date string) *cli.App {
	name := "otelgen"
	flags := getGlobalFlags()

	var v string
	if version == "" {
		v = "develop"
	} else {
		v = fmt.Sprintf("v%v-%v (%v) (ct)", version, commit, date)
	}

	app := &cli.App{
		Name:    name,
		Usage:   "A tool to generate synthetic OpenTelemetry logs, metrics and traces",
		Version: v,
		Flags:   flags,
		Commands: []*cli.Command{
			// genDiagnosticsCommand(),
			genLogsCommand(),
			genMetricsCommand(),
			genTracesCommand(),
		},
		Before: initLogger,
	}

	app.EnableBashCompletion = true

	return app
}
