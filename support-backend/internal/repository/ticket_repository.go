package repository

import (
	"context"
	"support-backend/internal/domain"
)

type TicketRepository interface {
	Create(ctx context.Context, t domain.Ticket) error
	GetByID(ctx context.Context, id string) (domain.Ticket, bool, error)
	List(ctx context.Context) ([]domain.Ticket, error)
	Update(ctx context.Context, t domain.Ticket) error
	Delete(ctx context.Context, id string) (bool, error)
}
