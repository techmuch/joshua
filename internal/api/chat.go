package api

import (
	"bd_bot/internal/ai"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ChatHandler struct {
	chatSvc *ai.ChatService
}

func NewChatHandler(svc *ai.ChatService) *ChatHandler {
	return &ChatHandler{chatSvc: svc}
}

type ChatRequest struct {
	Messages []ai.ChatMessage `json:"messages"`
}

func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// We can't use json.NewDecoder directly on Body if we want to support streaming properly 
	// because we need to set headers BEFORE reading the body if it takes time? 
	// No, we read request body first (fast), then start streaming response.
	
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Messages) == 0 {
		http.Error(w, "No messages provided", http.StatusBadRequest)
		return
	}

	// Set SSE Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	err := h.chatSvc.ChatStream(req.Messages, func(chunk string) error {
		// Wrap content in JSON object for frontend convenience
		// payload: { "content": "..." }
		// SSE format: data: <payload>\n\n
		payloadBytes, _ := json.Marshal(map[string]string{"content": chunk})
		fmt.Fprintf(w, "data: %s\n\n", payloadBytes)
		flusher.Flush()
		return nil
	})

	if err != nil {
		// Log error, but we can't send a 500 if we already started streaming headers.
		// In a production app, we might send a special error event.
		slog.Error("Chat stream failed", "error", err)
		
		// Optional: Send error event
		errPayload, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "data: %s\n\n", errPayload)
		flusher.Flush()
	}
}