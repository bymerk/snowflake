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

.PHONY: test
test:
	go test -v ./...

.PHONY: start-snowflake
start-snowflake:
	go run cmd/snowflake.go

.PHONY: docker-build
docker-build:
	@docker build --no-cache -t snowflake -f build/Dockerfile .
	docker tag snowflake bymerk/snowflake:latest

.PHONY: docker-push
docker-push:
	docker push bymerk/snowflake:latest
