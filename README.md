# HaruMQ

Minimal message broker in Go (MVP).

## Quick Start

### 1. Build

```
go build -o broker ./cmd/broker
go build -o producer ./cmd/producer
go build -o consumer ./cmd/consumer
```

### 2. Run Broker

```
./broker
```

### 3. Produce a Message

```
./producer --topic=orders --message="New order #123"
```

### 4. Consume Messages

```
./consumer --topic=orders --offset=0
```

#### Optional flags
- `--addr` (default: http://localhost:9092) — Broker address
- `--max` (consumer, default: 10) — Max messages to fetch

## Configuration

Edit `config.yaml` to set data directory and port:

```yaml
data_dir: "./data"
port: 9092
```

## How it works
- Broker exposes `/produce` (POST) and `/consume` (GET) endpoints.
- Producer CLI sends messages to a topic.
- Consumer CLI fetches messages from a topic and offset.
- Messages are stored in append-only log files per topic in `data/`.
