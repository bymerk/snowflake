package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr string
	GRPCAddr string
	NodeID   int64
	Log      bool
}

func LoadConfig() (*Config, error) {
	config := &Config{
		HTTPAddr: "0.0.0.0:8080",
		GRPCAddr: "0.0.0.0:5051",
		NodeID:   0,
	}

	if httpAddr, exists := os.LookupEnv("HTTP_ADDR"); exists {
		config.HTTPAddr = httpAddr
	}

	if grpcAddr, exists := os.LookupEnv("GRPC_ADDR"); exists {
		config.GRPCAddr = grpcAddr
	}

	if nodeIDEnv, exists := os.LookupEnv("NODE_ID"); exists {
		nodeID, err := strconv.Atoi(nodeIDEnv)
		if err != nil {
			return nil, err
		}
		config.NodeID = int64(nodeID)
	}

	return config, nil
}
