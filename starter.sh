set -e

echo "Pulling Docker images..."
docker pull postgres:16
docker pull redis:latest
docker pull golang:1.25.5-trixie

echo "Starting services with Docker Compose..."
docker compose up -d

echo "Installing golang-migrate CLI..."
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Ensure GOPATH/bin is available
export PATH="$PATH:$(go env GOPATH)/bin"

echo "🗄️  Running database migrations..."
migrate \
  -source file://migrations \
  -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" \
  up

echo "Setup completed successfully."
