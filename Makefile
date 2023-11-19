generate:
	protoc -I. --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/fish_catcher.proto

clean:
	rm proto/*.pb.go;

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

test:
	go test -cover -race ./...