FROM golang:1.24-alpine
RUN apk add --no-cache git curl
RUN go install github.com/air-verse/air@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["air"]