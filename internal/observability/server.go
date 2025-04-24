package observability

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	mux := http.NewServeMux()

	mux.HandleFunc("/live", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`OK`))
	})

	mux.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`OK`))
	})

	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to serve: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()
	ctxStop, cancelStop := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelStop()
	if err := srv.Shutdown(ctxStop); err != nil {
		log.Printf("failed to shutdown: %v \n", err)
	}

	return nil
}
