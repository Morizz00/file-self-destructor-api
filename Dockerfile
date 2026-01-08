FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary and static files
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/download.html .
COPY --from=builder /app/styles.css .
COPY --from=builder /app/script.js .

EXPOSE 8080

ENV PORT=8080

CMD ["./main"]

