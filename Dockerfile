FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1

COPY go.mod ./
COPY vendor/ ./vendor/
COPY . .

RUN go build -mod=vendor -o app ./cmd/main.go

# -------- Run Stage --------
FROM docker.io/library/alpine:latest

# SQLite runtime library
RUN apk add --no-cache sqlite-libs ca-certificates

RUN mkdir -p /app/data && chmod 777 /app/data

RUN adduser -D -u 1000 data && \
    chown -R data:data /app

USER data
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/utils ./utils

EXPOSE 8080

CMD ["./app"]