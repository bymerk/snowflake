package bootstrap

import (
	"context"

	"github.com/bymerk/snowflake/internal/config"
	grpcSF "github.com/bymerk/snowflake/internal/grpc"
	"github.com/bymerk/snowflake/internal/grpc/handler"
	"github.com/bymerk/snowflake/internal/http"
	"github.com/bymerk/snowflake/internal/observability"
	"github.com/bymerk/snowflake/pkg/showflake"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg *config.Config

	grpcSF     *grpcSF.Server
	httpSF     *http.Server
	metricsSrv *observability.Server
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	var sfOpts []showflake.Option

	if !cfg.Epoch.IsZero() {
		sfOpts = append(sfOpts, showflake.WithEpoch(cfg.Epoch))
	}

	sf, err := showflake.NewSnowflake(cfg.ClusterID, cfg.NodeID, sfOpts...)
	if err != nil {
		return nil, err
	}

	gsf := grpcSF.NewServer(
		grpcSF.Config{
			Addr:    cfg.GRPCAddr,
			Metrics: cfg.MetricsEnabled(),
		},
		handler.NewHandler(sf),
	)
	hsf := http.NewServer(
		http.Config{
			Addr:    cfg.HTTPAddr,
			Metrics: cfg.MetricsEnabled(),
		},
		sf,
	)

	var metricsSrv *observability.Server
	if cfg.MetricsEnabled() {
		metricsSrv = observability.NewServer(cfg.MetricsAddr)
	}

	return &App{
		cfg:        cfg,
		grpcSF:     gsf,
		httpSF:     hsf,
		metricsSrv: metricsSrv,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return app.grpcSF.Run(ctx)
	})

	errGroup.Go(func() error {
		return app.httpSF.Run(ctx)
	})

	if app.cfg.MetricsAddr != "" {
		errGroup.Go(func() error {
			return app.metricsSrv.Run(ctx)
		})
	}

	err := errGroup.Wait()
	if err != nil {
		return err
	}

	return nil
}
