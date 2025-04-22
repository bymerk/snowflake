package http

import (
	"context"
	"errors"
	"log"
	httpStd "net/http"
	"strconv"
	"time"

	"github.com/bymerk/snowflake/pkg/showflake"
)

type Server struct {
	addr string
	sf   *showflake.Snowflake
}

func NewServer(addr string, sf *showflake.Snowflake) *Server {
	return &Server{addr: addr, sf: sf}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	mux := httpStd.NewServeMux()

	mux.HandleFunc("/generate-id", func(writer httpStd.ResponseWriter, request *httpStd.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(httpStd.StatusOK)
		_, _ = writer.Write([]byte(`{"id":"`))
		_, _ = writer.Write([]byte(strconv.FormatInt(s.sf.Generate(), 10)))
		_, _ = writer.Write([]byte(`"}`))
	})

	srv := httpStd.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, httpStd.ErrServerClosed) {
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
