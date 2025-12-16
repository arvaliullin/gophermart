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


.PHONY: helloworld
helloworld:
	- go run github.com/arvaliullin/gophermart/examples/helloworld
