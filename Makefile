IMG_NAME=luongtruong20201/bookmark_service

GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

COVERAGE_FOLDER=./coverage

IMG_TAG := $(or $(GIT_TAG),$(BRANCH),dev)
export IMG_TAG

COVERAGE_EXCLUDE=mocks|main.go|test|infrastructure
COVERAGE_THRESHOLD=50

.PHONY: run
run: swagger
	go run cmd/api/main.go

.PHONY: test
test:
	CGO_ENABLED=1 go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -d . -o docs

.PHONY: redis
redis:
	docker run --name redis -p 6379:6379 redis:latest

.PHONY: docker-build
docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

.PHONY: docker-release
docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

.PHONY: docker-login
docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

.PHONY: docker-test
docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build \
		--build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" \
		--target test \
		-t $(IMG_NAME):$(IMG_TAG) \
		--output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

.PHONY: generate-rsa-key
generate-rsa-key:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private.pem -out public.pem

.PHONY: migrate
migrate:
	go run cmd/migrate/main.go

.PHONY: generate
generate:
	go generate ./...