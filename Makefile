.PHONY: build cover start test test-integration

build:
	docker build -t stockinos/api .

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main cmd/server/*

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

migrate:
	go run cmd/migrate/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...

ngrok:
	ngrok http --region=us --hostname=api.stockinos.ngrok.io 8000

protogen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/protos/*
