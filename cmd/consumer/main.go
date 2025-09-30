package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	topic := flag.String("topic", "", "Topic name")
	partition := flag.Int("partition", 0, "Partition number")
	offset := flag.Int64("offset", 0, "Offset to start consuming from")
	max := flag.Int("max", 10, "Max messages to fetch")
	addr := flag.String("addr", "http://localhost:9092", "Broker address")
	flag.Parse()

	if *topic == "" {
		fmt.Println("Usage: consumer --topic=TOPIC [--partition=N] [--offset=N] [--max=M] [--addr=ADDR]")
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/consume?topic=%s&partition=%d&offset=%d&max=%d", *addr, *topic, *partition, *offset, *max)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("Broker error:", string(body))
		os.Exit(1)
	}
	var msgs []map[string]interface{}
	if err := json.Unmarshal(body, &msgs); err != nil {
		fmt.Println("Decode error:", err)
		os.Exit(1)
	}
	for _, m := range msgs {
		payloadBase64, _ := m["payload"].(string)
		payload, err := base64.StdEncoding.DecodeString(payloadBase64)
		if err != nil {
			payload = []byte(payloadBase64)
		}
		fmt.Printf("[partition=%v offset=%v] %s\n", *partition, m["offset"], string(payload))
	}
}
