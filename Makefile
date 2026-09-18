-include .env
export

MAIN_PACKAGE_PATH = ./cmd/server/main.go

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
    # prints lines that start with '##' and use ':' as separator
    # example ## target: usage
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
    # $ is special to Make, to pass a literal $ to the shell, you must escape it as $$
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
    # `test -z STRING` checks if STRING is zero length
    # `git status --porcelain` produces a machine-readable list of changes in the repository.
	@test -z "$(shell git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	mkdir -p tmp/
	go test -v -race -buildvcs -coverprofile=tmp/coverage.out ./...
	go tool cover -html=tmp/coverage.out

## upgradeable: list direct dependencies that have upgrades available
.PHONY: upgradeable
upgradeable:
	go run github.com/oligot/go-mod-upgrade@latest

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and modernize and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fix ./...
	go fmt ./...

## dev/build: generate templ go, files and bundle js and css, and build server
.PHONY: dev/build
dev/build:
	go tool templ generate
	npm run dev:build
	go build -o ./tmp/main ./cmd/server/main.go

## dev: native Go + Docker Postgres
.PHONY: dev
dev: dev/clean
	docker run -d \
	 --name dev_postgres \
	 -p $(DB_PORT):5432 \
	 -e POSTGRES_USER=$(DB_USER) \
	 -e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	 -e POSTGRES_DB=$(DB_NAME) \
	 -v ./scripts:/docker-entrypoint-initdb.d \
	 postgres:16 && \
	docker run -d --name dev_redis -p ${REDIS_PORT}:6379 redis:latest && \
	DEV="true" DB_HOST="localhost" REDIS_HOST="localhost" go tool air -c .air.toml

## docker/status: show running containers
.PHONY: docker/status
docker/status:
	@docker ps --format "table {{.Names}}\\t{{.Status}}\\t{{.Ports}}"

## dev/clean: stops and removes dev_postgres and dev_redis container
.PHONY: dev/clean
dev/clean:
	@echo "Stop and removing dev dev_postgres container..."
	@docker stop dev_postgres 2>/dev/null || true
	@docker rm dev_postgres 2>/dev/null || true
	@echo "Stop and removing dev dev_redis container..."
	@docker stop dev_redis 2>/dev/null || true
	@docker rm dev_redis 2>/dev/null || true

## install: installs dependencies
.PHONY: install
install:
	@npm install
	@go mod download

## clean: clean up the build binaries
.PHONY: clean
clean: confirm dev/clean
	@echo "Cleaning up..."
	@rm -f web/templates/**/*_templ.go
	@rm -f cmd/server/static/app.js
	@rm -f cmd/server/static/app.css
	@rm -rf node_modules

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## db/dump-schema: makes schema (includes indexes)
.PHONY: db/dump-schema
db/dump-schema:
	@docker exec dev_postgres pg_dump -U $(DB_USER) --schema-only $(DB_NAME) > scripts/01-schema.sql

## db/dump-seed: makes new seed data
.PHONY: db/dump-seed
db/dump-seed:
	@docker exec dev_postgres pg_dump -U $(DB_USER) --data-only --inserts $(DB_NAME) > scripts/02-seed.sql
