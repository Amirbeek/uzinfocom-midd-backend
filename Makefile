-include .env
-include .envrc

# `go install` binarylarni $GOPATH/bin ga qo'yadi, u odatda PATH da bo'lmaydi.
# Shu qatorsiz `migrate` va `swag` "command not found" beradi.
export PATH := $(shell go env GOPATH)/bin:$(PATH)

DB_ADDR         ?= postgres://orderUser:adminpassword@localhost:5434/Order?sslmode=disable
MIGRATIONS_PATH ?= ./cmd/migrate/migrations

.PHONY: run migrate-up migrate-down migration install_tools

run:
	@go run ./cmd/api

migrate-up:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose up

migrate-down:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose down 1

migration:
	@if [ -z "$(name)" ]; then echo "Usage: make migration name=<snake_case_name>"; exit 1; fi
	@migrate create -seq -ext sql -dir "$(MIGRATIONS_PATH)" "$(name)"

install_tools:
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
