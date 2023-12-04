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

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

auto-run-server:
	ls cmd/server/*.go | entr -r go run cmd/server/main.go

auto-run-client:
	ls cmd/client/*.go | entr -r go run cmd/client/main.go

docker.run: docker.mongo docker.rabbitmq docker.redis
docker.stop: docker.mongo.stop docker.rabbitmq.stop docker.redis.stop

docker.mongo:
	docker run --rm -d \
		--name cgapp-mongo \
		-p 27017:27017 \
		mongo

docker.mongo.stop:
	docker stop cgapp-mongo

docker.rabbitmq:
	docker run --rm -d \
        --name cgapp-rabbitmq \
        -p 5672:5672 \
        -p 15672:15672 \
        -e RABBITMQ_DEFAULT_USER=guest \
        -e RABBITMQ_DEFAULT_PASS=guest \
        rabbitmq:3-management

docker.rabbitmq.stop:
	docker stop cgapp-rabbitmq

docker.redis:
	docker run --rm -d \
        --name cgapp-redis \
        -p 6379:6379 \
        -e REDIS_PASSWORD=myredispassword \
        redis

docker.redis.stop:
	docker stop cgapp-redis

generate-proto:
	protoc -I. --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/user.proto

clean-proto:
	rm proto/*.pb.go;
