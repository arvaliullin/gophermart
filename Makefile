.PHONY: install-deps
install-deps:
	- go install go.uber.org/mock/mockgen@latest
	- go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: generate-mocks
generate-mocks:
	go generate ./...

.PHONY: fmt
fmt:
	- go fmt ./...

.PHONY: up
up:
	- docker compose -f deployments/docker-compose.yaml up --build -d

.PHONY: down
down:
	- docker compose -f deployments/docker-compose.yaml down -v

.PHONY: logs
logs:
	- docker compose -f deployments/docker-compose.yaml logs

.PHONY: ps
ps:
	- docker compose -f deployments/docker-compose.yaml ps

.PHONY: prune
prune: down
	- docker image prune -f
	- docker container prune -f
	- docker volume prune -f
	- docker network prune -f
	- docker system prune -a --volumes -f


.PHONY: run
run:
	- go run github.com/arvaliullin/gophermart/cmd/gophermart

.PHONY: test
test:
	go test ./... -cover

.PHONY: test-v
test-v:
	go test ./... -v -cover

.PHONY: test-short
test-short:
	go test ./... -short -cover

.PHONY: test-integration
test-integration:
	go test -tags=integration ./... -cover

.PHONY: test-integration-v
test-integration-v:
	go test -tags=integration ./... -v -cover

.PHONY: build
build:
	go build -o bin/gophermart ./cmd/gophermart

.PHONY: lint
lint:
	golangci-lint run

GOPHERMART_HOST ?= localhost
GOPHERMART_PORT ?= 8080
ACCRUAL_PORT ?= 8081
DATABASE_URI ?= postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable

.PHONY: autotest
autotest: build
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
migration:
ifndef name
	$(error Использование: make migration name=<имя_миграции>)
endif
	goose -dir migrations create $(name) go
