# HaruMQ - Minimal Message Broker

## Build

```
make build
```

## Run Broker

```
make run-broker
```

## Run Producer

```
make run-producer TOPIC=orders MESSAGE="New order #123"
```

## Run Consumer

```
make run-consumer TOPIC=orders OFFSET=0
```

## Clean

```
make clean
```
