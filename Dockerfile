FROM golang:1.22-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/out cmd/app/main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/out /app/out
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY  ./migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["/bin/sh", "-c", "goose up && /app/out"]
