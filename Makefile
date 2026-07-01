.PHONY: build test run sync-types prod-build prod-build-amd64 prod-export prod-import prod-deploy prod-up prod-down prod-logs registration-code

PROD_COMPOSE := docker compose -f docker-compose.prod.yml
PROD_RUNTIME_COMPOSE := docker compose -f docker-compose.prod.yml -f docker-compose.prod.runtime.yml
PROD_IMAGES := financialtracker-proxy:latest financialtracker-backend:latest
PROD_EXPORT_ARCHIVE := dist/ft-prod-images.tar.gz
VPS_HOST ?=
VPS_PATH ?=
VPS_USER ?= $(shell whoami)

sync-types:
	@cd app && go run scripts/sync_types.go

build: sync-types
	@cd app && go build -o ../bin/main main.go

run: build
	@./bin/main

prod-build:
	$(PROD_COMPOSE) build

prod-build-amd64:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 $(PROD_COMPOSE) build

prod-export:
	@mkdir -p dist
	docker save $(PROD_IMAGES) | gzip > $(PROD_EXPORT_ARCHIVE)
	@echo "Exported $(PROD_EXPORT_ARCHIVE)"

prod-import:
	gunzip -c $(PROD_EXPORT_ARCHIVE) | docker load

prod-deploy:
	@test -n "$(VPS_HOST)" || (echo "Set VPS_HOST (e.g. user@your.vps)" && exit 1)
	@test -n "$(VPS_PATH)" || (echo "Set VPS_PATH (e.g. /home/user/FinancialTracker)" && exit 1)
	$(MAKE) prod-build-amd64
	$(MAKE) prod-export
	scp $(PROD_EXPORT_ARCHIVE) $(VPS_HOST):/tmp/ft-prod-images.tar.gz
	ssh $(VPS_HOST) 'gunzip -c /tmp/ft-prod-images.tar.gz | docker load'
	ssh $(VPS_HOST) 'cd $(VPS_PATH) && docker compose -f docker-compose.prod.yml -f docker-compose.prod.runtime.yml up -d'
	@echo "Deployed to $(VPS_HOST):$(VPS_PATH)"

prod-up:
	$(PROD_COMPOSE) up -d

prod-up-runtime:
	$(PROD_RUNTIME_COMPOSE) up -d

prod-down:
	$(PROD_COMPOSE) down

prod-logs:
	$(PROD_COMPOSE) logs -f

REGISTRATION_COMPOSE ?= $(PROD_COMPOSE)

registration-code:
	@if $(REGISTRATION_COMPOSE) ps --status running -q backend 2>/dev/null | grep -q .; then \
		$(REGISTRATION_COMPOSE) exec -T backend /app/main gen-registration-code; \
	else \
		$(REGISTRATION_COMPOSE) run --rm --no-deps -T backend /app/main gen-registration-code; \
	fi
