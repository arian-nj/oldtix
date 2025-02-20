updb:
	sudo docker compose up -d
	dbmate up

build:
	dbmate migrate
	sqlc generate

run_main:
	go run ./cmd/api/.
