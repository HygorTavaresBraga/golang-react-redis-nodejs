package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Event struct {
	Type      string    `json:"type"`
	TicketID  string    `json:"ticketId"`
	Timestamp time.Time `json:"timestamp"`
}

type Bus interface {
	Publish(ctx context.Context, e Event) error
	Subscribe(ctx context.Context) (<-chan Event, func(), error)
}

type RedisBus struct {
	rdb     *redis.Client
	channel string
}

func NewRedisBus(rdb *redis.Client, channel string) *RedisBus {
	return &RedisBus{rdb: rdb, channel: channel}
}

func (b *RedisBus) Publish(ctx context.Context, e Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return b.rdb.Publish(ctx, b.channel, payload).Err()
}

func (b *RedisBus) Subscribe(ctx context.Context) (<-chan Event, func(), error) {
	sub := b.rdb.Subscribe(ctx, b.channel)

	// garante que a inscrição está ativa (boa prática)
	if _, err := sub.Receive(ctx); err != nil {
		_ = sub.Close()
		return nil, nil, err
	}

	out := make(chan Event, 64)
	ch := sub.Channel()

	done := make(chan struct{})

	go func() {
		defer close(out)
		defer close(done)
		for msg := range ch {
			var e Event
			if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
				// ignora mensagens inválidas
				continue
			}
			select {
			case out <- e:
			case <-ctx.Done():
				return
			}
		}
	}()

	unsubscribe := func() {
		_ = sub.Close()
		<-done
	}

	return out, unsubscribe, nil
}
