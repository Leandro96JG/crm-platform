FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
COPY migrations/ ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/resources ./resources
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

ENV MIGRATIONS_PATH=file://migrations
ENV APP_PROPERTY_FILE=resources/application.properties

CMD ["./server"]
