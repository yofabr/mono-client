FROM golang:1.25.5-trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /mono-client ./cmd/main.go

FROM scratch
COPY --from=builder /mono-client /mono-client
EXPOSE 8080
CMD ["/mono-client"]
