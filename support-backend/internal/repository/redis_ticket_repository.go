package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"support-backend/internal/domain"

	"github.com/redis/go-redis/v9"
)

type RedisTicketRepository struct {
	rdb *redis.Client
}

func NewRedisTicketRepository(rdb *redis.Client) *RedisTicketRepository {
	return &RedisTicketRepository{rdb: rdb}
}

// Chaves usadas no Redis:
// ticket:<id>        -> JSON do ticket
// tickets:index      -> SET de ids (para listagem simples)
func ticketKey(id string) string { return "ticket:" + id }

const ticketsIndexKey = "tickets:index"

func (r *RedisTicketRepository) Create(ctx context.Context, t domain.Ticket) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	pipe := r.rdb.TxPipeline()
	pipe.Set(ctx, ticketKey(t.ID), b, 0)
	pipe.SAdd(ctx, ticketsIndexKey, t.ID)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisTicketRepository) GetByID(ctx context.Context, id string) (domain.Ticket, bool, error) {
	val, err := r.rdb.Get(ctx, ticketKey(id)).Result()
	if err == redis.Nil {
		return domain.Ticket{}, false, nil
	}
	if err != nil {
		return domain.Ticket{}, false, err
	}

	var t domain.Ticket
	if err := json.Unmarshal([]byte(val), &t); err != nil {
		return domain.Ticket{}, false, fmt.Errorf("corrupt ticket json: %w", err)
	}
	return t, true, nil
}

func (r *RedisTicketRepository) List(ctx context.Context) ([]domain.Ticket, error) {
	ids, err := r.rdb.SMembers(ctx, ticketsIndexKey).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []domain.Ticket{}, nil
	}

	out := make([]domain.Ticket, 0, len(ids))
	for _, id := range ids {
		t, ok, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, t)
		}
	}
	return out, nil
}

func (r *RedisTicketRepository) Update(ctx context.Context, t domain.Ticket) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	// garante que está no índice
	pipe := r.rdb.TxPipeline()
	pipe.Set(ctx, ticketKey(t.ID), b, 0)
	pipe.SAdd(ctx, ticketsIndexKey, t.ID)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisTicketRepository) Delete(ctx context.Context, id string) (bool, error) {
	// DEL retorna quantas chaves foram removidas
	n, err := r.rdb.Del(ctx, ticketKey(id)).Result()
	if err != nil {
		return false, err
	}
	// remove do índice (mesmo se não existir, não dá erro)
	if err := r.rdb.SRem(ctx, ticketsIndexKey, id).Err(); err != nil {
		return false, err
	}
	return n > 0, nil
}
