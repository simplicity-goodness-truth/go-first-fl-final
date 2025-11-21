# Stage 1: Builder stage
FROM golang:1.24.0-alpine AS builder

WORKDIR /app 

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tracker

# Stage 2: Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /tracker .

COPY web ./web

ENV TODO_DBFILE=scheduler.db
ENV TODO_PORT=7540
ENV TODO_PASSWORD=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6InVzZXIxIn0.NdVbWc1NZbtrS7IA4-3W-Tf8fg9pnk53vNHx9G1dU4A

EXPOSE ${TODO_PORT}

CMD ["./tracker"]

