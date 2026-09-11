-include .env
-include .envrc

# `go install` binarylarni $GOPATH/bin ga qo'yadi, u odatda PATH da bo'lmaydi.
# Shu qatorsiz `migrate` va `swag` "command not found" beradi.
export PATH := $(shell go env GOPATH)/bin:$(PATH)

GOBIN   := $(shell go env GOPATH)/bin
SWAG    := $(GOBIN)/swag
MIGRATE := $(GOBIN)/migrate

DB_ADDR         ?= postgres://orderUser:adminpassword@localhost:5434/Order?sslmode=disable
MIGRATIONS_PATH ?= ./cmd/migrate/migrations

.PHONY: run docs migrate-up migrate-down migrate-reset migration install_tools

docs:
	@$(SWAG) init -g cmd/api/main.go -o docs

run: docs
	@go run ./cmd/api


migrate-reset:
	@$(MIGRATE) -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" down -all
	@$(MIGRATE) -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose up
migrate-up:
	@$(MIGRATE) -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose up

migrate-down:
	@$(MIGRATE) -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose down 1

migration:
	@if [ -z "$(name)" ]; then echo "Usage: make migration name=<snake_case_name>"; exit 1; fi
	@$(MIGRATE) create -seq -ext sql -dir "$(MIGRATIONS_PATH)" "$(name)"

install_tools:
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
