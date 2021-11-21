gen:
	protoc --proto_path=proto --go_out=plugins=grpc:pb proto/*.proto
clean:
	rm -rf pb/*
run:
	go run main.go
test:
	go test -cover -race ./...