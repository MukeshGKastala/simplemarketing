EXECUTABLE_NAME := marketing-api

build: tidy
	@CGO_ENABLED=0 go build -o bin/$(EXECUTABLE_NAME) .

run: build
	@bin/$(EXECUTABLE_NAME)

tidy:
	@go fmt ./...
	@go mod tidy -v

test:
	@go test -v -count=1 -coverprofile=coverage.out ./...

coverage-func: test
	@go tool cover -func=coverage.out

coverage-html: test
	@go tool cover -html=coverage.out

generate:
	@go generate ./...
	@sqlc generate

clean:
	@rm -rf coverage.out bin/
	@go clean
	@docker compose down

compose:
	@docker compose up --build

.PHONY: build run tidy test coverage-func coverage-html generate clean compose