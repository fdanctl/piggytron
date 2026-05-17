-include .env
export

MAIN_PACKAGE_PATH = ./cmd/server/main.go

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

## audit: run quality control checks
.PHONY: audit
audit:
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## upgradeable: list direct dependencies that have upgrades available
.PHONY: upgradeable
upgradeable:
	go run github.com/oligot/go-mod-upgrade@latest

## tidy: tidy modfiles and modernize and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fix ./...
	go fmt ./...

## dev-build: generate templ go, files and bundle js and css, and build server
.PHONY: dev-build
dev-build:
	go tool templ generate
	npm run dev:build
	go build -o ./tmp/main ./cmd/server/main.go

## dev: native Go + Docker Postgres
.PHONY: dev
dev: clean-dev
	docker run -d \
	 --name postgres \
	 -p $(DB_PORT):5432 \
	 -e POSTGRES_USER=$(DB_USER) \
	 -e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	 -e POSTGRES_DB=$(DB_NAME) \
	 -v ./scripts:/docker-entrypoint-initdb.d \
	 postgres:16-alpine && \
	docker run -d --name redis -p ${REDIS_PORT}:6379 redis:latest && \
	DEV="true" go tool air -c .air.toml

## templ-watch: watch for templ files
.PHONY: templ-watch
templ-watch:
	DEV="true" go tool templ generate -watch -cmd="go run $(MAIN_PACKAGE_PATH)"

## status: show running containers
.PHONY: status
status:
	@docker ps --format "table {{.Names}}\\t{{.Status}}\\t{{.Ports}}"

## clean-dev: stops and removes dev postgres container
.PHONY: clean-dev
clean-dev:
	@echo "Stop and removing dev postgres container..."
	@docker stop postgres 2>/dev/null || true
	@docker rm postgres 2>/dev/null || true
	@echo "Stop and removing dev redis container..."
	@docker stop redis 2>/dev/null || true
	@docker rm redis 2>/dev/null || true

## clean: clean up the build binaries
.PHONY: clean
clean: confirm clean-dev
	@echo "Cleaning up..."
	@rm -f web/templates/**/*_templ.go
	@rm -f web/static/assets/app.js
	@rm -f web/static/assets/app.css
