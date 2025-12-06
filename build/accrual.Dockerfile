FROM debian:trixie-slim

WORKDIR /app

COPY cmd/accrual/accrual_linux_amd64 /app/accrual

CMD ["/app/accrual"]
