package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func startBroker() (*exec.Cmd, error) {
	cmd := exec.Command("go", "run", "./cmd/broker/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, cmd.Start()
}

func TestProduceConsume(t *testing.T) {
	// Start broker
	cmd, err := startBroker()
	if err != nil {
		t.Fatalf("failed to start broker: %v", err)
	}
	defer cmd.Process.Kill()
	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Produce message
	produceBody := map[string]interface{}{
		"topic":     "testtopic",
		"payload":   "hello world",
		"partition": 0,
	}
	b, _ := json.Marshal(produceBody)
	resp, err := http.Post("http://localhost:9092/produce", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("produce failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("produce status: %v", resp.Status)
	}
	resp.Body.Close()

	// Consume message
	url := "http://localhost:9092/consume?topic=testtopic&partition=0&offset=0&max=1"
	resp, err = http.Get(url)
	if err != nil {
		t.Fatalf("consume failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("consume status: %v", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	var msgs []map[string]interface{}
	if err := json.Unmarshal(body, &msgs); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(msgs) == 0 || msgs[0]["payload"] != "hello world" {
		t.Fatalf("unexpected consume result: %v", msgs)
	}
}
