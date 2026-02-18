package service

import (
	"context"
	"errors" //definição de erros customizados
	"log"
	"strings" //manipulação de strings (trim, toUpper, etc)
	"time"    //manipulação de tempo e datas

	"support-backend/internal/domain" //importa o pacote domain (definição da struct Ticket)
	"support-backend/internal/events"
	"support-backend/internal/repository"

	"github.com/google/uuid" //geração de UUIDs únicos
)

// Definição de erros customizados para o serviço de tickets
var (
	ErrNotFound      = errors.New("Recurso não encontrado")
	ErrInvalidInput  = errors.New("Campo de entrada inválido")
	ErrInvalidStatus = errors.New("Status inválido")
)

// Estruturas de entrada para criação e atualização de tickets (DTOs)
type CreateTicketInput struct { //DTO para criar ticket (ainda tem tags JSON porque handler usa essa struct para decodificar JSON)
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTicketInput struct { //DTO para atualizar ticket (no UPDATE, pode ser adicionado ponteiros (*string) nos campos opcionais de preenchimento na atualização)
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type TicketService struct {
	repo repository.TicketRepository
	bus  events.Bus
}

func NewTicketService(repo repository.TicketRepository, bus events.Bus) *TicketService {
	return &TicketService{repo: repo, bus: bus}
}

func (s *TicketService) publishEvent(e events.Event) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.bus.Publish(ctx, e); err != nil {
		log.Printf("ticket event publish failed: %v", err)
	}
}

func (s *TicketService) Create(ctx context.Context, input CreateTicketInput) (domain.Ticket, error) {
	title := strings.TrimSpace(input.Title)
	desc := strings.TrimSpace(input.Description)

	if len(title) < 3 || len(title) > 120 {
		return domain.Ticket{}, ErrInvalidInput
	}
	if len(desc) < 3 || len(desc) > 2000 {
		return domain.Ticket{}, ErrInvalidInput
	}

	now := time.Now().UTC()
	t := domain.Ticket{
		ID:          uuid.NewString(),
		Title:       title,
		Description: desc,
		Status:      domain.TicketOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return domain.Ticket{}, err
	}

	s.publishEvent(events.Event{
		Type:      "ticket.created",
		TicketID:  t.ID,
		Timestamp: time.Now().UTC(),
	})

	return t, nil
}

func (s *TicketService) GetByID(ctx context.Context, id string) (domain.Ticket, error) {
	t, ok, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Ticket{}, err
	}
	if !ok {
		return domain.Ticket{}, ErrNotFound
	}
	return t, nil
}

func (s *TicketService) List(ctx context.Context) ([]domain.Ticket, error) {
	tickets, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (s *TicketService) Update(ctx context.Context, id string, input UpdateTicketInput) (domain.Ticket, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Ticket{}, ErrInvalidInput
	}

	// 1) Carrega o ticket atual
	t, ok, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Ticket{}, err
	}
	if !ok {
		return domain.Ticket{}, ErrNotFound
	}

	// 2) Aplica alterações (somente se o campo veio)
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if len(title) < 3 || len(title) > 120 {
			return domain.Ticket{}, ErrInvalidInput
		}
		t.Title = title
	}

	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		if len(desc) < 3 || len(desc) > 2000 {
			return domain.Ticket{}, ErrInvalidInput
		}
		t.Description = desc
	}

	if input.Status != nil {
		normalized := domain.TicketStatus(strings.ToUpper(strings.TrimSpace(*input.Status)))
		switch normalized {
		case domain.TicketOpen, domain.TicketPending, domain.TicketClosed:
			t.Status = normalized
		default:
			return domain.Ticket{}, ErrInvalidStatus
		}
	}

	// 3) Atualiza timestamp e persiste
	t.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, t); err != nil {
		return domain.Ticket{}, err
	}

	s.publishEvent(events.Event{
		Type:      "ticket.updated",
		TicketID:  t.ID,
		Timestamp: time.Now().UTC(),
	})

	return t, nil
}

func (s *TicketService) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrInvalidInput
	}

	deleted, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}

	s.publishEvent(events.Event{
		Type:      "ticket.deleted",
		TicketID:  id,
		Timestamp: time.Now().UTC(),
	})

	return nil
}
