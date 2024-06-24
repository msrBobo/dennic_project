CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

	./scripts/genproto.sh	migrate -path migrations -database $(DB_URL) -verbose up

DB_URL := "postgres://postgres:mubina2007@localhost:5432/dennic_session_service?sslmode=disable"

migrate-up:


migrate-down:

proto-gen:
	./scripts/genproto.sh

migrate-force:
	migrate -path migrations -database $(DB_URL) -verbose force 1

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_session_table

pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge
