FROM golang:1.25.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/gophermart github.com/arvaliullin/gophermart/cmd/gophermart

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/local/bin/gophermart /usr/local/bin/gophermart

CMD ["gophermart"]
