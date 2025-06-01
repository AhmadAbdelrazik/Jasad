include .env

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: run/mqtt
run/mqtt:
	go run ./cmd/mqtt

.PHONY: run/http
run/http:
	go run ./cmd/http

.PHONY: db/psql
db/psql:
	psql ${JASAD_DB_DSN}

.PHONY: db/migrations/new
db/migration/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${JASAD_DB_DSN} up

.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rolling down migrations...'
	migrate -path ./migrations -database ${JASAD_DB_DSN} down

