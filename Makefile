INTERNAL_PKG = ./internal/...
ALL_PKG = ./internal/... ./pkg/...
COV_DIR = scripts/result_cover

include .env
export

start:
	mkdir -p -m 777 logs/app
	docker compose down && docker compose up -d

start-rebuild:
	mkdir -p -m 777 logs/app
	docker compose down && docker compose up -d --build

run-local:
	mkdir -p -m 777 logs/app
	POSTGRES_HOST=localhost go run ./cmd/main.go -config=./cmd/config/config.yml

stop:
	docker compose down

rm-docker-volume:
	docker compose down --volumes

cov:
	mkdir -p ${COV_DIR}
	go test -race -coverpkg=${INTERNAL_PKG} -coverprofile ${COV_DIR}/cover.out ${INTERNAL_PKG}; cat ${COV_DIR}/cover.out | fgrep -v "test.go" | fgrep -v "docs" | fgrep -v "mock" | fgrep -v "config" > ${COV_DIR}/cover2.out
	go tool cover -func ${COV_DIR}/cover2.out
	go tool cover -html ${COV_DIR}/cover2.out -o ${COV_DIR}/coverage.html

test:
	go test ./...