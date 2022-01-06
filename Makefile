install:
	go mod download
gen:
	protoc --proto_path=proto --go_out=plugins=grpc:pb proto/**/*.proto
clean:
	rm -rf pb/*
server:
	go run cmd/server/main.go --port 8080
client:
	go run cmd/client/main.go --address 0.0.0.0:8080
test:
	go test -cover -race ./...

.PHONY: gen clean server client test