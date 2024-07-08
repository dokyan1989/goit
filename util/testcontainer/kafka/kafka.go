package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

type Container struct {
	container *kafka.KafkaContainer
	address   string
}

func Run(ctx context.Context) (*Container, error) {
	kafkaC, err := kafka.RunContainer(ctx,
		testcontainers.WithImage("confluentinc/cp-kafka:latest"),
		kafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	containerHost, err := kafkaC.Host(ctx)
	if err != nil {
		log.Printf("failed to get container host: %s\n", err)
		return nil, err
	}

	mappedPort, err := kafkaC.MappedPort(ctx, "9093/tcp")
	if err != nil {
		log.Printf("failed to get container port: %s\n", err)
		return nil, err
	}

	return &Container{
		container: kafkaC,
		address:   fmt.Sprintf("%s:%s", containerHost, mappedPort.Port()),
	}, nil
}

func (c *Container) GetAddress() string {
	return c.address
}

func (c *Container) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}
