FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x && go mod verify

COPY . ./

RUN go build -v -o main ./cmd/server
# RUN GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" -o /bin/server ./cmd/server

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main /app/main

CMD ["/app/main"]
