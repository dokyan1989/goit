package redis

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestContainer(t *testing.T) {
	ctx := context.Background()

	c, err := Run(ctx)
	if err != nil {
		t.Errorf("RunPostgresContainer() error = %v", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.GetConnString(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err = rdb.Ping(ctx).Err()
	if err != nil {
		t.Errorf("rdb.Ping() error = %v", err)
		return
	}

	err = c.Terminate(ctx)
	if err != nil {
		t.Errorf("c.Terminate() error = %v", err)
		return
	}
}
