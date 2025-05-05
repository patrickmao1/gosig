FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./node ./cmd/gosig_test

FROM alpine:latest

WORKDIR /app

RUN mkdir runtime

COPY --from=builder /app/node .
COPY --from=builder /app/test_config.yaml ./gosig_config.yaml

CMD ./node 2>&1 | tee ./runtime/app.log