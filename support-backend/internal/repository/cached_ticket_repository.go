package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"support-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CachedTicketRepository struct {
	next *RedisTicketRepository
	rdb  *redis.Client
	ttl  time.Duration

	lockTTL  time.Duration
	waitMax  time.Duration
	waitStep time.Duration
}

func NewCachedTicketRepository(next *RedisTicketRepository, rdb *redis.Client, ttl time.Duration) *CachedTicketRepository {
	return &CachedTicketRepository{
		next: next,
		rdb:  rdb,
		ttl:  ttl,

		lockTTL:  2 * time.Second,        // lock curto (evita thundering herd)
		waitMax:  800 * time.Millisecond, // quanto tempo no máximo outros requests esperam
		waitStep: 60 * time.Millisecond,  // intervalo entre tentativas
	}
}

func cacheKey(id string) string { return "ticket:cache:" + id }

// Create: escreve no storage e já popula cache
func (c *CachedTicketRepository) Create(ctx context.Context, t domain.Ticket) error {
	if err := c.next.Create(ctx, t); err != nil {
		return err
	}
	return c.setCache(ctx, t)
}

// GetByID: tenta cache primeiro; se miss, busca storage e popula cache
func (c *CachedTicketRepository) GetByID(ctx context.Context, id string) (domain.Ticket, bool, error) {
	// 1) tentar cache direto
	if t, ok, err := c.getFromCache(ctx, id); err != nil {
		return domain.Ticket{}, false, err
	} else if ok {
		return t, true, nil
	}

	// 2) tentar pegar lock
	token, locked, err := c.tryLock(ctx, id)
	if err != nil {
		return domain.Ticket{}, false, err
	}

	if locked {
		// Sou o "reconstrutor" do cache
		defer c.unlock(ctx, id, token)

		// double-check (alguém pode ter preenchido em microtempo)
		if t, ok, err := c.getFromCache(ctx, id); err != nil {
			return domain.Ticket{}, false, err
		} else if ok {
			return t, true, nil
		}

		// buscar do storage
		t, ok, err := c.next.GetByID(ctx, id)
		if err != nil {
			return domain.Ticket{}, false, err
		}
		if !ok {
			// opcional: cache negativo por poucos segundos para reduzir spam em IDs inexistentes
			// (vamos manter simples e não fazer agora)
			return domain.Ticket{}, false, nil
		}

		// preencher cache (best effort)
		_ = c.setCache(ctx, t)
		return t, true, nil
	}

	// 3) não peguei lock -> vou esperar um pouco e tentar ler cache de novo
	deadline := time.Now().Add(c.waitMax)
	for time.Now().Before(deadline) {
		time.Sleep(c.waitStep)

		if t, ok, err := c.getFromCache(ctx, id); err != nil {
			return domain.Ticket{}, false, err
		} else if ok {
			return t, true, nil
		}
	}

	// 4) fallback: se ninguém preencheu o cache a tempo, vai ao storage
	// (isso garante que o usuário não fica preso esperando)
	t, ok, err := c.next.GetByID(ctx, id)
	if err != nil {
		return domain.Ticket{}, false, err
	}
	if !ok {
		return domain.Ticket{}, false, nil
	}
	_ = c.setCache(ctx, t)
	return t, true, nil
}

// List: deixa sem cache (por enquanto)
func (c *CachedTicketRepository) List(ctx context.Context) ([]domain.Ticket, error) {
	return c.next.List(ctx)
}

// Update: escreve storage e reescreve cache
func (c *CachedTicketRepository) Update(ctx context.Context, t domain.Ticket) error {
	if err := c.next.Update(ctx, t); err != nil {
		return err
	}
	return c.setCache(ctx, t)
}

// Delete: remove storage e remove cache
func (c *CachedTicketRepository) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := c.next.Delete(ctx, id)
	if err != nil {
		return false, err
	}
	_ = c.rdb.Del(ctx, cacheKey(id)).Err()
	return deleted, nil
}

func (c *CachedTicketRepository) setCache(ctx context.Context, t domain.Ticket) error {
	if c.ttl <= 0 {
		// TTL 0 => cache desativado
		return nil
	}

	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	if err := c.rdb.Set(ctx, cacheKey(t.ID), b, c.ttl).Err(); err != nil {
		return fmt.Errorf("cache set failed: %w", err)
	}
	return nil
}

func lockKey(id string) string { return "ticket:lock:" + id }

func (c *CachedTicketRepository) tryLock(ctx context.Context, id string) (string, bool, error) {
	token := uuid.NewString()
	ok, err := c.rdb.SetNX(ctx, lockKey(id), token, c.lockTTL).Result()
	if err != nil {
		return "", false, err
	}
	return token, ok, nil
}

// unlock seguro: só remove se o token bater (evita deletar lock de outro)
func (c *CachedTicketRepository) unlock(ctx context.Context, id string, token string) {
	const lua = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end
`
	_ = c.rdb.Eval(ctx, lua, []string{lockKey(id)}, token).Err()
}

func (c *CachedTicketRepository) getFromCache(ctx context.Context, id string) (domain.Ticket, bool, error) {
	val, err := c.rdb.Get(ctx, cacheKey(id)).Result()
	if err == redis.Nil {
		return domain.Ticket{}, false, nil
	}
	if err != nil {
		return domain.Ticket{}, false, err
	}

	var t domain.Ticket
	if err := json.Unmarshal([]byte(val), &t); err != nil {
		// cache corrompido -> remove
		_ = c.rdb.Del(ctx, cacheKey(id)).Err()
		return domain.Ticket{}, false, nil
	}
	return t, true, nil
}
