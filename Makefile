.PHONY: mocks
mocks:
	mockgen -source=internal/quote/domain/quote.go -destination=internal/quote/domain/mock/quote.go

.PHONY: deps
deps:
	go mod tidy
	go mod vendor
	
.PHONY: api
api:
	go run cmd/api/main.go

.PHONY: tests
tests:
	go test ./...