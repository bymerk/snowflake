package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bymerk/snowflake/internal/config"
	grpcSF "github.com/bymerk/snowflake/internal/grpc"
	"github.com/bymerk/snowflake/internal/grpc/handler"
	"github.com/bymerk/snowflake/internal/http"
	"github.com/bymerk/snowflake/pkg/showflake"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	sf, err := showflake.NewSnowflake(cfg.NodeID)
	if err != nil {
		log.Fatalf("Error creating snowflake generator: %v", err)
	}

	gsf := grpcSF.NewServer(cfg.GRPCAddr, handler.NewHandler(sf))
	hsf := http.NewServer(cfg.HTTPAddr, sf)

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return gsf.Run(ctx)
	})

	errGroup.Go(func() error {
		return hsf.Run(ctx)
	})

	err = errGroup.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
