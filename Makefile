run:
	go run ./cmd/main.go

tidy:
	go mod tidy

dev:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env and configure it."; \
		exit 1; \
	fi
	docker build -t mono-client:dev -f Dockerfile.dev .
	docker run --rm -it --env-file .env -v "$$(pwd):/app" mono-client:dev

dev-build:
	docker build -t mono-client:dev -f Dockerfile.dev .

.PHONY: up
up:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env and configure it."; \
		exit 1; \
	fi
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: logs
logs:
	docker-compose logs -f