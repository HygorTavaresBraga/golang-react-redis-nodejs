package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"

	"support-backend/internal/config"
	"support-backend/internal/events"
	"support-backend/internal/handler"
	"support-backend/internal/middleware"
	"support-backend/internal/repository"
	"support-backend/internal/service"
)

func main() {
	cfg := config.Load()

	// 1) Redis client (único modo do projeto)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	// 2) Pub/Sub bus (streaming)
	bus := events.NewRedisBus(rdb, "tickets:events")
	streamHandler := handler.NewTicketStreamHandler(bus)

	// 3) Repository (storage real + cache TTL)
	storageRepo := repository.NewRedisTicketRepository(rdb)
	ttl := time.Duration(cfg.TicketCacheTTLSeconds) * time.Second
	repo := repository.NewCachedTicketRepository(storageRepo, rdb, ttl)

	// 4) Service + Handlers
	ticketSvc := service.NewTicketService(repo, bus)
	ticketHandler := handler.NewTicketHandler(ticketSvc)

	// 5) Routes
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/api/v1/tickets/stream", streamHandler.Stream)

	mux.HandleFunc("/api/v1/tickets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ticketHandler.List(w, r)
		case http.MethodPost:
			ticketHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/tickets/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ticketHandler.GetByID(w, r)
		case http.MethodPatch:
			ticketHandler.Update(w, r)
		case http.MethodDelete:
			ticketHandler.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 6) Middlewares
	var h http.Handler = mux
	h = middleware.MaxBodyBytes(1<<20, h)
	h = middleware.SecurityHeaders(h)
	h = middleware.RequestID(h)
	h = middleware.Logger(h)

	// 7) HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 0,
		IdleTimeout:  60 * time.Second,
	}

	// start
	go func() {
		log.Printf("server running on :%s (env=%s)", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
