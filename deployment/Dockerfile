FROM golang:1.24.1-alpine

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY .. .
COPY ../templates ./templates

RUN go build -o mail-service cmd/main.go

EXPOSE 8080

CMD ["./mail-service"]
