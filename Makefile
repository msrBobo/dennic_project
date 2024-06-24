-include .env
export

CURRENT_DIR=$(shell pwd)
APP=Booking_service
CMD_DIR=./cmd

.DEFAULT_GOAL = build
POSTGRES_USER = postgres
POSTGRES_PASSWORD = 1234
POSTGRES_HOST = localhost
POSTGRES_PORT = 5544
POSTGRES_DATABASE = dennic

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

.PHONY: create-migration
create-migration:
	migrate create -dir migrations -ext sql -seq $(name)_table

# proto-gen
.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

DB_URL := "postgres://postgres:20030505@localhost:5432/dennic?sslmode=disable"

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down

.PHONY: migrate-file
migrate-file:
	migrate create -ext sql -dir migrations/ -seq alter_appointments_table

# run test
test:
	go test -v -cover -race ./internal/...

# -------------- for deploy --------------
build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge
