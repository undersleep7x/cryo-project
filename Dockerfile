FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod go.sum ./
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
    postgresql-client \
    make
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.2
RUN go mod tidy && go build -o ./tmp/main .
RUN rm -f tmp/main
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]