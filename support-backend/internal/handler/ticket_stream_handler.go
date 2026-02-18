package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"support-backend/internal/events"
)

type TicketStreamHandler struct {
	bus events.Bus
}

func NewTicketStreamHandler(bus events.Bus) *TicketStreamHandler {
	return &TicketStreamHandler{bus: bus}
}

func (h *TicketStreamHandler) Stream(w http.ResponseWriter, r *http.Request) {
	// SSE headers obrigatórios
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Alguns proxies só “liberam” streaming com isso:
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	eventsCh, unsubscribe, err := h.bus.Subscribe(ctx)
	if err != nil {
		http.Error(w, "failed to subscribe", http.StatusInternalServerError)
		return
	}
	defer unsubscribe()

	// Opcional: evento inicial (handshake)
	_ = writeSSE(w, "connected", map[string]any{
		"ts": time.Now().UTC().Format(time.RFC3339),
	})
	flusher.Flush()

	// Loop de eventos
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			// Comentario SSE para manter a conexao viva (evita timeouts/proxies fechando).
			_, _ = fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
		case e, ok := <-eventsCh:
			if !ok {
				return
			}
			_ = writeSSE(w, "ticket", e)
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, event string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Formato SSE:
	// event: <nome>
	// data: <json>
	// \n
	_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, string(b))
	return err
}
