package handler

import (
	"context"

	"github.com/bymerk/snowflake/internal/grpc/gen"
	"github.com/bymerk/snowflake/pkg/showflake"
)

type Handler struct {
	sf *showflake.Snowflake
	gen.UnimplementedSnowflakeServiceServer
}

func NewHandler(sf *showflake.Snowflake) *Handler {
	return &Handler{
		sf: sf,
	}
}

func (h *Handler) GenerateID(_ context.Context, _ *gen.GenerateIDRequest) (*gen.GenerateIDResponse, error) {
	return &gen.GenerateIDResponse{Id: h.sf.Generate()}, nil
}
