package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/krzko/otelgen/internal/metrics"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
)

var generateMetricsUpDownCounterCommand = &cli.Command{
	Name:        "up-down-counter",
	Usage:       "generate metrics of type up down counter",
	Description: "UpDownCounter demonstrates how to measure numbers that can go up and down",
	Aliases:     []string{"udc"},
	Hidden:      true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "temporality",
			Usage: "Temporality defines the window that an aggregation was calculated over, one of: delta, cumulative",
			Value: "delta",
		},
	},
	Action: func(c *cli.Context) error {
		return generateMetricsUpDownCounterAction(c)
	},
	Before: func(c *cli.Context) error {
		fmt.Println("DEPRECATION WARNING: The 'counter' command is deprecated and will be removed in a future version.")
		fmt.Println("Please use the 'guage' command instead.")
		fmt.Println("Example: otelgen metrics guage")
		fmt.Println()
		return nil
	},
}

func generateMetricsUpDownCounterAction(c *cli.Context) error {
	var err error

	if c.String("otel-exporter-otlp-endpoint") == "" {
		return errors.New("'otel-exporter-otlp-endpoint' must be set")
	}

	metricsCfg := &metrics.Config{
		TotalDuration: time.Duration(c.Duration("duration")),
		Endpoint:      c.String("otel-exporter-otlp-endpoint"),
		Rate:          c.Int64("rate"),
		ServiceName:   c.String("service-name"),
	}

	configureLogging(c)

	grpcExpOpt, httpExpOpt := getExporterOptions(c, metricsCfg)

	ctx := context.Background()

	exp, err := createExporter(ctx, c, grpcExpOpt, httpExpOpt)
	if err != nil {
		logger.Error("failed to obtain OTLP exporter", zap.Error(err))
		return err
	}
	defer shutdownExporter(exp)

	logger.Info("Starting metrics generation")

	reader := metric.NewPeriodicReader(
		exp,
		metric.WithInterval(time.Duration(metricsCfg.Rate)),
	)

	provider := createMeterProvider(reader, metricsCfg)

	metrics.SimulateUpDownCounter(provider, metricsCfg, logger)

	return nil
}
