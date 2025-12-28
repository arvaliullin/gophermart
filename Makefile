.DEFAULT_GOAL := help

.PHONY: help
help: ## Показать список доступных команд
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: install-deps
install-deps: ## Установить зависимости (mockgen, goose)
	- go install go.uber.org/mock/mockgen@latest
	- go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: generate-mocks
generate-mocks: ## Сгенерировать моки
	go generate ./...

.PHONY: fmt
fmt: ## Форматировать код
	- go fmt ./...

.PHONY: up
up: ## Запустить все сервисы через Docker Compose
	- docker compose -f deployments/docker-compose.yaml up --build -d

.PHONY: down
down: ## Остановить и удалить контейнеры Docker Compose
	- docker compose -f deployments/docker-compose.yaml down -v

.PHONY: logs
logs: ## Показать логи Docker Compose
	- docker compose -f deployments/docker-compose.yaml logs

.PHONY: ps
ps: ## Показать статус контейнеров Docker Compose
	- docker compose -f deployments/docker-compose.yaml ps

.PHONY: prune
prune: down ## Полная очистка Docker (образы, контейнеры, volumes)
	- docker image prune -f
	- docker container prune -f
	- docker volume prune -f
	- docker network prune -f
	- docker system prune -a --volumes -f


.PHONY: run
run: ## Запустить приложение локально
	- go run github.com/arvaliullin/gophermart/cmd/gophermart

.PHONY: test
test: ## Запустить все тесты с покрытием
	go test ./... -cover

.PHONY: test-v
test-v: ## Запустить все тесты с подробным выводом
	go test ./... -v -cover

.PHONY: test-short
test-short: ## Запустить короткие тесты
	go test ./... -short -cover

.PHONY: test-integration
test-integration: ## Запустить интеграционные тесты
	go test -tags=integration ./... -cover

.PHONY: test-integration-v
test-integration-v: ## Запустить интеграционные тесты с подробным выводом
	go test -tags=integration ./... -v -cover

.PHONY: build
build: ## Собрать бинарник приложения
	go build -o bin/gophermart ./cmd/gophermart

.PHONY: lint
lint: ## Запустить линтер
	golangci-lint run

GOPHERMART_HOST ?= localhost
GOPHERMART_PORT ?= 8080
ACCRUAL_PORT ?= 8081
DATABASE_URI ?= postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable

.PHONY: autotest
autotest: build ## Запустить автотесты
	./cmd/gophermarttest/gophermarttest \
		-test.v \
		-gophermart-binary-path=./bin/gophermart \
		-gophermart-host=$(GOPHERMART_HOST) \
		-gophermart-port=$(GOPHERMART_PORT) \
		-gophermart-database-uri="$(DATABASE_URI)" \
		-accrual-binary-path=./cmd/accrual/accrual_linux_amd64 \
		-accrual-host=$(GOPHERMART_HOST) \
		-accrual-port=$(ACCRUAL_PORT) \
		-accrual-database-uri="$(DATABASE_URI)"

.PHONY: migration
migration: ## Создать новую миграцию (использование: make migration name=имя_миграции)
ifndef name
	$(error Использование: make migration name=<имя_миграции>)
endif
	goose -dir migrations create $(name) go
