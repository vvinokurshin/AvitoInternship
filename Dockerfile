FROM golang:1.20

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
RUN go build cmd/main.go

EXPOSE 8001