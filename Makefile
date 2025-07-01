APP_NAME := key-server
SERVER_PORT := 8080

.PHONY: build run test clean docker-build docker-run fmt help

help: ## Show available targets
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

build: ## Build the key-server binary
	go build -o $(APP_NAME) ./cmd/$(APP_NAME)

run: build ## Run the key-server locally
	./$(APP_NAME) --srv-port=$(SERVER_PORT)

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	go clean
	rm -f $(APP_NAME)

docker-build: ## Build Docker image
	docker build -t $(APP_NAME) .

docker-run: docker-build ## Run Docker container
	docker run -p $(SERVER_PORT):$(SERVER_PORT) $(APP_NAME)

fmt: ## Format Go code
	go fmt ./...

