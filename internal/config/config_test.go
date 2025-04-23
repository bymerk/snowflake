package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bymerk/snowflake/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	cfg := config.Config{
		HTTPAddr:    "127.0.0.1:8080",
		GRPCAddr:    "127.0.0.1:5051",
		MetricsAddr: "127.0.0.1:9000",
		ClusterID:   1,
		NodeID:      11,
		Epoch:       time.UnixMicro(1609459200000),
	}

	envMap := map[string]string{
		"HTTP_ADDR":    cfg.HTTPAddr,
		"GRPC_ADDR":    cfg.GRPCAddr,
		"METRICS_ADDR": cfg.MetricsAddr,
		"CLUSTER_ID":   fmt.Sprintf("%d", cfg.ClusterID),
		"NODE_ID":      fmt.Sprintf("%d", cfg.NodeID),
		"EPOCH":        fmt.Sprintf("%d", cfg.Epoch.UnixMilli()),
	}

	for k, v := range envMap {
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}

	c, err := config.LoadConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, *c)
}
