COVERAGE_EXCLUDE=mocks|main.go|test
COVERAGE_THRESHOLD=50

.PHONY: run
run: swagger
	go run cmd/api/main.go

.PHONY: test
test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -d . -o docs

.PHONY: redis
redis:
	docker run  --name redis -p 6379:6379 redis:latest

.PHONY: docker-build
docker-build:
	docker build -f Dockerfile .