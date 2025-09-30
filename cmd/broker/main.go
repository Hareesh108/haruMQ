package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/yourname/harumq/internal/api"
	"github.com/yourname/harumq/internal/broker"
	"github.com/yourname/harumq/internal/storage"
)

func main() {
	cfg, err := broker.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logic, err := storage.NewLog(cfg.DataDir)
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}
	server := &api.Server{Log: logic}
	http.HandleFunc("/produce", server.Produce)
	http.HandleFunc("/consume", server.Consume)
	addr := ":" + strconv.Itoa(cfg.Port)
	log.Printf("Broker listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
