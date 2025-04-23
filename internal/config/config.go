package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr    string
	GRPCAddr    string
	MetricsAddr string

	ClusterID int64
	NodeID    int64
	Epoch     time.Time // milliseconds

}

func LoadConfig() (*Config, error) {
	config := &Config{
		HTTPAddr: "0.0.0.0:8080",
		GRPCAddr: "0.0.0.0:5051",
	}

	if httpAddr, exists := os.LookupEnv("HTTP_ADDR"); exists {
		config.HTTPAddr = httpAddr
	}

	if grpcAddr, exists := os.LookupEnv("GRPC_ADDR"); exists {
		config.GRPCAddr = grpcAddr
	}

	if clusterEnv, exists := os.LookupEnv("CLUSTER_ID"); exists {
		clusterID, err := strconv.Atoi(clusterEnv)
		if err != nil {
			return nil, err
		}
		config.ClusterID = int64(clusterID)
	}

	if nodeIDEnv, exists := os.LookupEnv("NODE_ID"); exists {
		nodeID, err := strconv.Atoi(nodeIDEnv)
		if err != nil {
			return nil, err
		}
		config.NodeID = int64(nodeID)
	}

	if epochEnv, exists := os.LookupEnv("EPOCH"); exists {
		epoch, err := strconv.Atoi(epochEnv)
		if err != nil {
			return nil, err
		}
		config.Epoch = time.UnixMilli(int64(epoch))
	}

	if metricsAddrEnv, exists := os.LookupEnv("METRICS_ADDR"); exists {
		config.MetricsAddr = metricsAddrEnv
	}

	return config, nil
}
