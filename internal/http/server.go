package http

import (
	"context"
	"errors"
	"log"
	httpStd "net/http"
	"strconv"
	"time"

	"github.com/bymerk/snowflake/internal/http/middleware"
	"github.com/bymerk/snowflake/pkg/showflake"
)

type Server struct {
	server *httpStd.Server
}

type Config struct {
	Addr    string
	Metrics bool
}

func NewServer(cfg Config, sf *showflake.Snowflake) *Server {
	var handler httpStd.Handler = getMux(sf)

	if cfg.Metrics {
		handler = middleware.MetricsMiddleware(handler)
	}

	srv := &httpStd.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	return &Server{server: srv}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, httpStd.ErrServerClosed) {
			log.Fatalf("failed to serve: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()
	ctxStop, cancelStop := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelStop()
	if err := s.server.Shutdown(ctxStop); err != nil {
		log.Printf("failed to shutdown: %v \n", err)
	}

	return nil
}

func getMux(sf *showflake.Snowflake) *httpStd.ServeMux {
	mux := httpStd.NewServeMux()
	mux.HandleFunc("/generate-id", func(writer httpStd.ResponseWriter, request *httpStd.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(httpStd.StatusOK)
		_, _ = writer.Write([]byte(`{"id":"`))
		_, _ = writer.Write([]byte(strconv.FormatInt(sf.Generate(), 10)))
		_, _ = writer.Write([]byte(`"}`))
	})
	return mux
}
