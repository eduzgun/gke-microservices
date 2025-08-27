# Generate protobuf code
proto-gen:
	mkdir -p ./internal/proto/gen/go/session
	protoc \
		--go_out=./internal/proto/gen/go/session \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/proto/gen/go/session \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./internal/proto \
		./internal/proto/session.proto

# Clean generated files
clean-proto:
	rm -rf ./internal/proto/gen/

# Generate and run session service
session-service: proto-gen
	go run cmd/session-service/main.go

# Generate and run main app
app: proto-gen
	go run cmd/app/main.go