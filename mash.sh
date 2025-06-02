#!/bin/bash

set -e

# --- Configurations --- 
export DATABASE_URL=postgres://pgadmin:pgpass@127.0.0.1:5432/game?sslmode=disable
export RELEASE_MODE=dev 
export BASE_URL=https://hokmran.filelord.ir
export RELEASE_MODE=dev
export SECRET_KEY=sdasfdgewfef
export PROJECT_VERSION=1.0.0 \

COMMAND=$1

generate_sqlc() {
	echo "Generating..."
	cd ./backend/
	sqlc generate
	cd ..
	echo "OK"	
}

if [ -z "COMMAND" ]; then
	echo "i want more commands"
fi

case $COMMAND in
	run_core)
		generate_sqlc
		echo "Starting Core API"
		cd ./backend/
		CORE_HTTP_PORT=4444 \
			go run ./cmd/core/.

		echo "The End"
		;;
	
	run_hokm4)	
		generate_sql
		echo "Starting Hokm Bacbkend"
		cd ./backend/
		HOKM4_HTTP_PORT=4445 \
		  go run ./cmd/hokm4/.
		echo "The End"
		;;

	run_bot)
		generate_sql
		echo "Starting Hokm Bacbkend"
		cd ./backend/
		WEBAPP_HTTP_PORT=3000 \
		  go run ./cmd/telegram/.
		echo "The End"
		;;
	
	gensqlc)
		generate_sqlc
	;;
	caddy)
		echo "Starting Caddy"
		sudo caddy run --config ./Caddyfile
		;;

	*)
		echo "command '$COMMAND' is Unknown"
esac

exit 0
