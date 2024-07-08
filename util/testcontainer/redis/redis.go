package redis

import (
	"context"
	"log"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type Container struct {
	container  *redis.RedisContainer
	connString string
}

func Run(ctx context.Context) (*Container, error) {
	redisC, err := redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:latest"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		log.Printf("failed to start container: %s\n", err)
		return nil, err
	}

	connString, err := redisC.ConnectionString(ctx)
	if err != nil {
		log.Printf("failed to get connection string: %s\n", err)
		return nil, err
	}

	return &Container{
		container:  redisC,
		connString: strings.TrimPrefix(connString, "redis://"),
	}, nil
}

func (c *Container) GetConnString() string {
	return c.connString
}

func (c *Container) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}
