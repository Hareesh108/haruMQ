package storage

import "time"

type Message struct {
	Topic   string    `json:"topic"`
	Offset  int64     `json:"offset"`
	Payload []byte    `json:"payload"`
	Ts      time.Time `json:"ts"`
}
