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
./producer --topic=orders --message="New order #123" --partition=0
```

### 4. Consume Messages

```
./consumer --topic=orders --partition=0 --offset=0
```

#### Optional flags
- `--partition` (default: 0) — Partition number (for both producer and consumer)
- `--addr` (default: http://localhost:9092) — Broker address
- `--max` (consumer, default: 10) — Max messages to fetch

## Configuration

Edit `config.yaml` to set data directory and port:

```yaml
data_dir: "./data"
port: 9092
```

## How it works

## Feature Comparison: Apache Kafka vs HaruMQ MVP

| Feature                        | Apache Kafka         | HaruMQ MVP (Current)         |
|--------------------------------|----------------------|------------------------------|
| Multi-topic support            | Yes                  | Yes                          |
| Partitioned topics             | Yes                  | No                           |
| Replication                    | Yes                  | No                           |
| Consumer groups                | Yes                  | No                           |
| Message durability             | Yes                  | Basic (append-only file)     |
| Log compaction                 | Yes                  | No                           |
| Retention policies             | Yes                  | No                           |
| High throughput                | Yes                  | Basic                        |
| Horizontal scalability         | Yes                  | No (single broker)           |
| Schema registry                | Yes (optional)       | No                           |
| REST API                       | Optional (Confluent) | Yes                          |
| CLI tools                      | Yes                  | Yes (basic)                  |
| Offset management              | Yes                  | Manual (via CLI param)       |
| Authentication/Authorization   | Yes                  | No                           |
| Monitoring/metrics             | Yes                  | No                           |
| Transactions                   | Yes                  | No                           |
| Web UI                         | Optional             | No                           |


### What HaruMQ MVP supports now
- Multiple topics
- Partitioned topics (per topic, multiple partitions)
- Append-only log storage (per topic/partition)
- Basic durability (messages survive restarts)
- Simple REST API for produce/consume (with partition support)
- Basic CLI for producer and consumer (with partition support)

### Not yet implemented (for production readiness)
- Partitioned topics
- Replication
- Consumer groups
- Log compaction and retention
- Security (auth, TLS)
- Monitoring/metrics
- Schema registry
- Transactions
- Web UI
