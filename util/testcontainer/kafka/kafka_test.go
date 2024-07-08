package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	c, err := Run(ctx)
	if err != nil {
		t.Errorf("failed to Run() error = %v", err)
		return
	}

	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", c.GetAddress(), topic, partition)
	if err != nil {
		t.Errorf("failed to dial leader: %v", err)
		return
	}

	// to produce messages
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		t.Errorf("failed to write messages: %v", err)
		return
	}

	// to consume messages
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// https://stackoverflow.com/a/78515123
	batch := conn.ReadBatch(1e1, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		t.Errorf("failed to close batch: %v", err)
		return
	}

	if err := conn.Close(); err != nil {
		t.Errorf("failed to close writer: %v", err)
		return
	}

	err = c.Terminate(ctx)
	if err != nil {
		t.Errorf("failed to Terminate() error = %v", err)
		return
	}
}
