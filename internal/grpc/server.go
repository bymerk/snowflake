package grpc

import (
	"context"
	"log"
	"net"

	"github.com/bymerk/snowflake/internal/grpc/gen"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

type Server struct {
	addr    string
	handler gen.SnowflakeServiceServer
}

func NewServer(addr string, handler gen.SnowflakeServiceServer) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)

	grpcServer := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	}...)
	gen.RegisterSnowflakeServiceServer(grpcServer, s.handler)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
