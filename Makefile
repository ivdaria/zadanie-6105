DB_DSN="postgres://postgres:postgres@localhost:5432/zadanie"
MIGRATIONS_DIR=migrations

install-deps: install-lint
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrations-up:
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_DSN} up

migrations-down:
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_DSN} down

migrations-status:
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_DSN} status

install-oapi-gen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

gen-spec:
	go generate -run="oapi-codegen" -tags="tools" -x ./...

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: install-lint
	golangci-lint run
