.PHONY: clean critic security lint test build run

APP_NAME = chatapp
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

generate:
	protoc -I. --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/weather.proto

docker.run: docker.mongo

docker.mongo:
	docker run --rm -d \
		--name cgapp-mongo \
		-p 27017:27017 \
		mongo

clean-proto:
	rm proto/*.pb.go;

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go