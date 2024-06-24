-include .env
export

CURRENT_DIR=$(shell pwd)
APP=dennic_user_service
CMD_DIR=./cmd

.DEFAULT_GOAL = build
POSTGRES_USER = postgres
POSTGRES_PASSWORD = 123
POSTGRES_HOST = localhost
POSTGRES_PORT = 5432
POSTGRES_DATABASE = dennic

# build for current os
.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

# migrate
.PHONY: migrate-up
migrate-up:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

# migrate
.PHONY: migrate-down
migrate-down:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable down

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_table_users

# proto
.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

# git submodule init 	
.PHONY: pull-proto
pull-proto:
	git submodule update --init --recursive

# go generate	
.PHONY: go-gen
go-gen:
	go generate ./...

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

.PHONY: pull-proto-module
pull-proto-module:
	git submodule update --init --recursive

.PHONY: update-proto-module
update-proto-module:
	git submodule update --remote --merge


#SERVER_HOST = dennic_api_gateway:
#SERVER_PORT = 9050
#
#REDIS_HOST = redisdb
#REDIS_PORT = 6379
#
#BOOKING_SERVICE_GRPC_HOST = dennic_booking_service
#BOOKING_SERVICE_GRPC_PORT = 9090
#
#HEALTHCARE_SERVICE_GRPC_HOST = dennic_healthcare_service
#HEALTHCARE_SERVICE_GRPC_PORT = 9080
#
#SESSION_SERVICE_GRPC_HOST = dennic_session_service
#SESSION_SERVICE_GRPC_PORT = 9060
#
#USER_SERVICE_GRPC_HOST = dennic_user_service
#USER_SERVICE_GRPC_PORT = 9070
#
#OTLP_COLLECTOR_HOST = otlp-collector
#
#MINIO_SERVICE_ENDPOINT = minio:9000
#
#POSTGRES_USER = postgres
#POSTGRES_PASSWORD = 20030505
#POSTGRES_HOST = postgresdb
#POSTGRES_PORT = 5432
#POSTGRES_DATABASE = dennic