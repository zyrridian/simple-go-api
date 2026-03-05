# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY . .

# Initialize go mod if it doesn't exist, then build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server .

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api-server .
EXPOSE 8081
CMD ["./api-server"]