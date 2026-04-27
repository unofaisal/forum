
FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY vendor/ ./vendor/
COPY . .

RUN go build -mod=vendor -o app

# -------- Run Stage --------
FROM docker.io/library/alpine:latest

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/app .

# Copy required project folders
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/utils ./utils

EXPOSE 8080

CMD ["./app"]