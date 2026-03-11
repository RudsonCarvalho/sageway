.PHONY: all build test lint bench proto docker dev clean

GOFLAGS := -race -count=1
COVERAGE_OUT := coverage.out
COVERAGE_THRESHOLD := 90

all: lint test build

build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/orchestrator ./cmd/orchestrator

test:
	go test ./... $(GOFLAGS) -coverprofile=$(COVERAGE_OUT)
	go tool cover -func=$(COVERAGE_OUT)

lint:
	golangci-lint run ./...

bench:
	go test ./... -bench=. -benchmem -benchtime=10s -run='^$$'

proto:
	protoc --go_out=. --go-grpc_out=. proto/l2/v1/l2_messages.proto

docker:
	docker build -f deployments/Dockerfile.gateway -t aasg-gateway:dev .
	docker build -f deployments/Dockerfile.orchestrator -t aasg-orchestrator:dev .

dev:
	docker-compose -f deployments/docker-compose.yml up -d etcd elasticsearch opa
	go run ./cmd/orchestrator &
	go run ./cmd/gateway

clean:
	rm -f bin/gateway bin/orchestrator $(COVERAGE_OUT)

coverage-check:
	@go tool cover -func=$(COVERAGE_OUT) | awk '/total:/{if ($$3+0 < $(COVERAGE_THRESHOLD)) {print "FAIL: coverage " $$3 " < $(COVERAGE_THRESHOLD)%"; exit 1}}'
