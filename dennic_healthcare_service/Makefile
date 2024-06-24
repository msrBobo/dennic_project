-include .env
export

CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go


proto-gen:
	./scripts/genproto.sh

DB_URL := "postgres://postgres:123@localhost:5432/dennic_healthcare_service?sslmode=disable"

migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down

migrate-force:
	migrate -path migrations -database $(DB_URL) -verbose force 1

migrate-file:
	migrate create -ext sql -dir migrations/ -seq alter_specialization_table

pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge

run:
	go run ${CMD_DIR}/app/main.go

#---------LOCAL---------

#OTLP_COLLECTOR_HOST = 0.0.0.0

CMD_DIR=./cmd
POSTGRES_USER = postgres
POSTGRES_PASSWORD = 123
POSTGRES_HOST = localhost
POSTGRES_PORT = 5432
POSTGRES_DATABASE = dennic

#---------DOCKER---------
# OTLP_COLLECTOR_HOST = otlp-collector

# POSTGRES_USER = postgres
# POSTGRES_PASSWORD = 20030505
# POSTGRES_HOST = postgresdb
# POSTGRES_PORT = 5432
# POSTGRES_DATABASE = dennic

# {
#   "error": "2 UNKNOWN: ERROR: missing FROM-clause entry for table \"dwh\" (SQLSTATE 42P01)"
# }

