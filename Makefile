.PHONY: cover start test test-integration

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main cmd/server/*

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...
