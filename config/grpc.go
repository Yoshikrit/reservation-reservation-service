package config

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcConfig struct {
	InventoryGrpcAddress string `env:"INVENTORY_GRPC_ADDRESS,required"`
}

type GrpcConns struct {
	Inventory *grpc.ClientConn
}

func NewGrpcConns(cfg GrpcConfig) (*GrpcConns, error) {
	inventoryConn, err := grpc.NewClient(cfg.InventoryGrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}
	return &GrpcConns{Inventory: inventoryConn}, nil
}

func (c *GrpcConns) Close() {
	c.Inventory.Close()
}
