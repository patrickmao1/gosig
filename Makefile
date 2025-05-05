.PHONY: proto
proto:
	protoc -I proto/ \
		--go_out=./ \
		--go-grpc_out=require_unimplemented_servers=false:./ \
		proto/*.proto

test_network_up:
	docker compose -f ./tests/network/compose.yaml --project-directory . up --build -d

down:
	docker compose -f tests/gosig/compose.yaml --project-directory . down

up: down
	docker compose -f tests/gosig/compose.yaml --project-directory . up --build -d
	docker compose -f tests/gosig/compose.yaml --project-directory . logs -f

prof:
	go tool pprof -http=localhost:8000 "http://localhost:6060/debug/pprof/profile?seconds=58"

config:
	go run ./cmd/config_gen/main.go

install:
	go build -o ./gosig ./cmd/gosig_test

start:
	./gosig -configPath="./config.yaml"