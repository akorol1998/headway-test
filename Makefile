SERVICE := payment-api

.PHONY: install
install:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: go-build
go-build:
	go build -o app ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY:test
test:
	go test -v payment-api/internal/services/payment

.PHONY:integration-test
integration-test: up-test
	go test -v payment-api/integration_tests
	$(MAKE) clean-test

# Docker part

.PHONY: build
build:
	@docker build -t $(SERVICE) --no-cache .

# Brings up the database for integration tests
.PHONY: up-test
up-test:
	@echo "Launching testing database"
	docker-compose -f docker-compose-tests.yml up -d

.PHONY: up
up:
	docker-compose up $(name)

.PHONY: down
down:
	docker-compose stop

.PHONY: clean-test
clean-test:
	@echo "Stopping testing database"
	docker-compose -f docker-compose-tests.yml stop
