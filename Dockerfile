# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy source and resolve dependencies (generates go.sum if missing).
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o spotsync-api .

# ---- Run stage ----
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/spotsync-api .
EXPOSE 8080
CMD ["./spotsync-api"]
