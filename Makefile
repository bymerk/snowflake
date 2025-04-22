.PHONY: proto
proto:
	protoc \
	  --proto_path=api/grpc/snowflake \
      --go_out=paths=source_relative:internal/grpc/gen \
      --go-grpc_out=paths=source_relative:internal/grpc/gen \
      snowflake.proto

.PHONY: lint
lint:
	golangci-lint run

.PHONY: start-snowflake
start-snowflake:
	go run cmd/snowflake.go
