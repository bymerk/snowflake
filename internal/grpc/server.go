package grpc

import (
	"context"
	"log"
	"net"

	"github.com/bymerk/snowflake/internal/grpc/gen"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	addr   string
}

type Config struct {
	Addr    string
	Metrics bool
}

func NewServer(cfg Config, handler gen.SnowflakeServiceServer) *Server {
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	}

	if cfg.Metrics {
		options = append(options, grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	}

	grpcServer := grpc.NewServer(options...)
	gen.RegisterSnowflakeServiceServer(grpcServer, handler)

	return &Server{
		server: grpcServer,
		addr:   cfg.Addr,
	}
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()
	s.server.GracefulStop()

	return nil
}
