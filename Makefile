LOCAL_BIN := $(CURDIR)/bin
BUF_BUILD := $(LOCAL_BIN)/buf

.bin_deps: export GOBIN := $(LOCAL_BIN)
.bin_deps:
	$(info Installing binary dependencies...)
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.buf-generate:
	$(info run buf generate...)
	PATH="$(LOCAL_BIN):${PATH}" $(BUF_BUILD) generate

generate: .buf-generate .tidy

.tidy:
	go mod tidy

.PHONY: \
	.bin_deps

DB_URL = postgres://postgres:password@localhost:5432/postgres?sslmode=disable

run:
	go run cmd/main.go

up:
	COMPOSE_PROJECT_NAME=go-grpc-crud docker compose up --build -d

down:
	COMPOSE_PROJECT_NAME=go-grpc-crud docker compose down

down-v:
	COMPOSE_PROJECT_NAME=go-grpc-crud docker compose down -v

# Применить все миграции
migrate:
	migrate -source file://migrations -database "$(DB_URL)" up

# Применить только одну следующую миграцию
migrate-one:
	migrate -source file://migrations -database "$(DB_URL)" up 1

# Откатить последнюю миграцию
migrate-down:
	migrate -source file://migrations -database "$(DB_URL)" down 1

# Принудительно поставить версию (например force 0)
force:
	migrate -source file://migrations -database "$(DB_URL)" force 0