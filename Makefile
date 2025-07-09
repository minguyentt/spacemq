run-rdb-dev:
	@go run ./cmd/redis-test-server

run-rdb-test:
	@go test ./internal/redisdbq
