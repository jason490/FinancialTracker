.PHONY: build test run sync-types prod-build prod-up prod-down prod-logs

sync-types:
	@cd app && go run scripts/sync_types.go

build: sync-types
	@cd app && go build -o ../bin/main main.go

run: build
	@./bin/main

prod-build:
	docker compose -f docker-compose.prod.yml build

prod-up:
	docker compose -f docker-compose.prod.yml up -d

prod-down:
	docker compose -f docker-compose.prod.yml down

prod-logs:
	docker compose -f docker-compose.prod.yml logs -f
