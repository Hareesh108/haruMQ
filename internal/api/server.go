package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/yourname/harumq/internal/storage"
)

type Server struct {
	Log *storage.Log
}

func (s *Server) Produce(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Topic   string `json:"topic"`
		Payload string `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	msg := &storage.Message{
		Topic:   req.Topic,
		Payload: []byte(req.Payload),
		Ts:      time.Now(),
	}
	offset, err := s.Log.Append(req.Topic, msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"offset": offset})
}

func (s *Server) Consume(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	offsetStr := r.URL.Query().Get("offset")
	maxStr := r.URL.Query().Get("max")
	if topic == "" || offsetStr == "" {
		http.Error(w, "missing topic or offset", 400)
		return
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid offset", 400)
		return
	}
	max := 10
	if maxStr != "" {
		if m, err := strconv.Atoi(maxStr); err == nil {
			max = m
		}
	}
	msgs, err := s.Log.Read(topic, offset, max)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}
