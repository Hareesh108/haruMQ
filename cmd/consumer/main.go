package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"


func main() {
	topic := flag.String("topic", "", "Topic name")
	offset := flag.Int64("offset", 0, "Offset to start consuming from")
	max := flag.Int("max", 10, "Max messages to fetch")
	addr := flag.String("addr", "http://localhost:9092", "Broker address")
	flag.Parse()

	if *topic == "" {
		fmt.Println("Usage: consumer --topic=TOPIC [--offset=N] [--max=M] [--addr=ADDR]")
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/consume?topic=%s&offset=%d&max=%d", *addr, *topic, *offset, *max)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
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
		fmt.Printf("[offset=%v] %s\n", m["offset"], m["payload"])
	}
}
