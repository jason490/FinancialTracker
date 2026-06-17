.PHONY: build test run templ tailwind sync-types

sync-types:
	@cd app && go run scripts/sync_types.go

templ:
	@cd app && templ generate

tailwind:
	npx tailwindcss -i ./app/web/static/css/input.css -o ./app/web/static/css/tailwind.css --watch

build: templ sync-types
	@cd app && go build -o ../bin/main main.go

run: build
	@./bin/main
