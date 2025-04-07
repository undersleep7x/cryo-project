FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cosmtrek/air@v1.40.4
RUN apk add --no-cache \
    git \
    curl \
    build-base \
    gcc \
    musl-dev \
    libc-dev \
    binutils \
    libgcc \
    make
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]