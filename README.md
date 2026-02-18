# Support Tickets - Full Stack (Go + React)

Sistema simples de tickets com backend em Go e frontend em React.
Suporta CRUD de tickets e atualizacao em tempo real via Server-Sent Events (SSE).

## Estrutura do Projeto

- support-backend/ - API Go + Redis (persistencia) + stream SSE
- support-front-end/ - UI React (Vite)

## Stack de Tecnologias

Backend:
- Go 1.25.x
- net/http (servidor HTTP)
- Redis 7 (persistencia + cache + pub/sub)
- SSE (atualizacao em tempo real)

Frontend:
- React 19
- Vite 7
- React Router
- Ionicons (icones)

Infra:
- Docker Compose (Redis)

## Como Funciona (Fluxo Completo)

1) O frontend carrega a lista chamando GET /api/v1/tickets.
2) O frontend abre uma conexao SSE em /api/v1/tickets/stream.
3) Ao criar/editar/deletar, o backend publica um evento no Redis Pub/Sub.
4) O handler de SSE recebe do canal e transmite para todos os clientes conectados.
5) O frontend recebe o evento e faz refresh da lista.

## Arquitetura do Backend

Pecas principais:
- handlers: endpoints HTTP e stream SSE.
- service: validacoes e publicacao de eventos.
- repository: persistencia no Redis + cache opcional.
- events bus: Redis Pub/Sub para eventos em tempo real.

Fluxo principal:
- cmd/api/main.go conecta Redis, repositorios, servicos, handlers e rotas.
- internal/handler/ticket_handler.go expoe o CRUD.
- internal/handler/ticket_stream_handler.go faz o SSE.
- internal/service/ticket_service.go valida e publica eventos.
- internal/events/bus.go publica/assina o canal Redis.
- internal/repository/redis_ticket_repository.go persiste os tickets.
- internal/repository/cached_ticket_repository.go faz cache por id.

## Arquitetura do Frontend

Pecas principais:
- src/api/http.js: wrapper de fetch com tratamento de erro.
- src/api/tickets.js: chamadas da API de tickets.
- src/app/router.jsx: rotas.
- Paginas:
  - TicketsList: lista + SSE
  - TicketDetail: edicao
  - TicketCreate: criacao
- Componentes:
  - TicketForm: formulario reutilizavel
  - Layout: cabecalho e navegacao

## Endpoints da API

- GET /health
- GET /api/v1/tickets
- POST /api/v1/tickets
- GET /api/v1/tickets/{id}
- PATCH /api/v1/tickets/{id}
- DELETE /api/v1/tickets/{id}
- GET /api/v1/tickets/stream (SSE)

## Como Rodar Localmente

### 1) Subir o Redis

Dentro de support-backend/:

```
docker compose up -d
```

### 2) Subir o Backend

Dentro de support-backend/:

```
go run ./cmd/api
```

A API sobe em http://localhost:8080.

Variaveis de ambiente (opcional):
- PORT (padrao 8080)
- ENV (padrao dev)
- REDIS_ADDR (padrao localhost:6379)
- REDIS_PASSWORD (padrao vazio)
- TICKET_CACHE_TTL_SECONDS (padrao 30)

### 3) Subir o Frontend

Dentro de support-front-end/:

```
npm install
npm run dev
```

A UI sobe em http://localhost:5173 e faz proxy de /api para o backend.

## Observacoes

- SSE precisa de conexoes longas; o WriteTimeout do servidor foi desativado.
- Redis e usado para persistencia e para eventos via Pub/Sub.
- Status de ticket: OPEN, PENDING, CLOSED.
