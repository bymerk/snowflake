package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bymerk/snowflake/internal/bootstrap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(ctx)
	if err != nil {
		log.Println(err)
	}
}
