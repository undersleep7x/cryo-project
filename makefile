APP_NAME=cryowallet
DOCKER_IMAGE=cryowallet-dev
COVERAGE_THRESHOLD=80.0
TEST_PATHS=./internal/prices/...

.PHONY: ci docker-ci docker-build docker-up docker-down clean lint test coverage docker-refresh

# Lint inside container
lint:
	docker compose run --rm -e GOFLAGS="-buildvcs=false" app golangci-lint run ./...

# Run tests + coverage inside container
test:
	docker compose run --rm -e GOFLAGS="-buildvcs=false" app sh -c "\
		go test -v -coverprofile=coverage.out ${TEST_PATHS} && \
		go tool cover -func=coverage.out"

coverage:
	docker compose run --rm -e GOFLAGS="-buildvcs=false" app sh -c '\
		go test -v -coverprofile=coverage.out ${TEST_PATHS} && \
		go tool cover -func=coverage.out && \
		coverage=$$(go tool cover -func=coverage.out | grep total: | awk '\''{print $$3}'\'' | sed '\''s/%//'\''); \
		echo Parsed coverage value: $$coverage% && \
		awk -v cov=$$coverage -v thresh=$(COVERAGE_THRESHOLD) '\''BEGIN { exit (cov+0 < thresh) ? 1 : 0 }'\'' || \
		( echo "Coverage ($$coverage%) is below threshold ($(COVERAGE_THRESHOLD)%). Failing." && exit 1 )'

# Run everything together in container (your CI mimic)
docker-ci: docker-build docker-up lint test coverage docker-down

# Build container
docker-build:
	docker compose build --no-cache

# Start containers
docker-up:
	docker compose up

# Tear down containers
docker-down:
	docker compose down -v --remove-orphans

docker-refresh:
	docker compose down -v --remove-orphans
	docker compose build --no-cache
	docker compose up

# Cleanup
clean:
	rm -f coverage.out