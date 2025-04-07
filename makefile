APP_NAME=cryowallet
DOCKER_IMAGE=cryowallet-dev
COVERAGE_THRESHOLD=80.0
TEST_PATHS=./internal/prices/...

.PHONY: build lint test coverage run docker-up docker-down docker-reset docker-build clean docker-test-clean ci

lint:
	golangci-lint run ./...

test:
	go test -v -coverprofile=coverage.out ./...

coverage:
	go tool cover -html=coverage.out

coverage-check:
	@coverage=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	awk -v cov=$$coverage -v thresh=$(COVERAGE_THRESHOLD) 'BEGIN { exit (cov+0 < thresh) ? 1 : 0 }' || \
	( echo "Coverage ($$coverage%) is below threshold ($(COVERAGE_THRESHOLD)%). Failing." && exit 1 )

run:
	go run ./cmd/main.go

docker-build:
	docker build -t ${DOCKER_IMAGE}

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-test-clean:
	docker-compose down -v && docker-compose run --rm app sh -c "\
	go test -v -coverprofile=coverage.out ${TEST_PATHS} && \
	go tool cover -func=coverage.out && \
	coverage=\$$(go tool cover -func=coverage.out | grep total: | awk '{print \$$3}' | sed 's/%//'); \
	echo Parsed coverage value: \$${coverage}% && \
	awk -v cov=\$${coverage} -v thresh=$(COVERAGE_THRESHOLD) 'BEGIN { exit (cov+0 < thresh) ? 1 : 0 }' || \
	( echo 'Coverage (\$${coverage}%) is below threshold ($(COVERAGE_THRESHOLD)%). Failing.' && exit 1 )"

docker-reset:
	docker-compose down -v

clean:
	rm -f coverage.out

ci: clean lint docker-test-clean docker-reset