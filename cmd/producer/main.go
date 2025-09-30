package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	topic := flag.String("topic", "", "Topic name")
	message := flag.String("message", "", "Message payload")
	partition := flag.Int("partition", 0, "Partition number")
	addr := flag.String("addr", "http://localhost:9092", "Broker address")
	flag.Parse()

	if *topic == "" || *message == "" {
		fmt.Println("Usage: producer --topic=TOPIC --message=MSG [--partition=N] [--addr=ADDR]")
		os.Exit(1)
	}

	body, _ := json.Marshal(map[string]interface{}{
		"topic":     *topic,
		"payload":   *message,
		"partition": *partition,
	})
	resp, err := http.Post(*addr+"/produce", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Broker error:", resp.Status)
		os.Exit(1)
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Printf("Produced to topic '%s' partition %d at offset %v\n", *topic, *partition, res["offset"])
}
