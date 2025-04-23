protogen:
	protoc -I proto/ \
		--go_out=./ \
		--go-grpc_out=./ \
		proto/*.proto

test_network_up:
	docker compose -f ./tests/network/compose.yaml --project-directory . up --build -d

gosig_up:
	docker compose -f tests/gosig/compose.yaml --project-directory . up --build -d
	docker compose -f tests/gosig/compose.yaml --project-directory . logs -f

config:
	go run ./cmd/config_gen/main.go