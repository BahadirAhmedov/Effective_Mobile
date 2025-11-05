FROM golang:1.24.7-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY local.env ./local.env
COPY config/local.yaml ./local.yaml

COPY . .

RUN go build -o /bin/app ./cmd/data-aggregation
RUN go build -o /bin/migrator ./cmd/migrator

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/app /bin/app
COPY --from=builder /bin/migrator /bin/migrator
COPY --from=builder /app/local.yaml /app/local.yaml

CMD ["/bin/app"]
