.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: test
test:
	go test -v ./...