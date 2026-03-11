.PHONY: all build test lint bench proto docker dev clean

all: lint test build

# ── Build ──────────────────────────────────────────────────────────────────
build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/orchestrator ./cmd/orchestrator

# ── Test ───────────────────────────────────────────────────────────────────
test:
	go test ./... -race -count=1 -coverprofile=coverage.out -timeout 5m
	go tool cover -func=coverage.out

# ── Lint ───────────────────────────────────────────────────────────────────
lint:
	golangci-lint run ./...

# ── Benchmarks ─────────────────────────────────────────────────────────────
bench:
	go test ./... -bench=. -benchmem -benchtime=10s -run='^$$'

# ── Protobuf ───────────────────────────────────────────────────────────────
proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/l2/v1/l2_messages.proto

# ── Docker ─────────────────────────────────────────────────────────────────
docker:
	docker build -f deployments/Dockerfile.gateway -t sageway-gateway:dev .
	docker build -f deployments/Dockerfile.orchestrator -t sageway-orchestrator:dev .

# ── Dev (local stack) ──────────────────────────────────────────────────────
dev:
	docker compose -f deployments/docker-compose.yml up -d
	go run ./cmd/orchestrator &
	go run ./cmd/gateway

# ── Clean ──────────────────────────────────────────────────────────────────
clean:
	rm -rf bin/ coverage.out
	docker compose -f deployments/docker-compose.yml down
