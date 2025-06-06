#!/bin/bash

set -e

# --- Configurations --- 
export DATABASE_URL=postgres://pgadmin:pgpass@127.0.0.1:5432/game?sslmode=disable
export RELEASE_MODE=dev 
export BASE_URL=https://hokmran.filelord.ir
export RELEASE_MODE=dev
export SECRET_KEY=sdasfdgewfef
export PROJECT_VERSION=1.0.0
export ANDROID_HOME="/nonexistent"

GODOT_VERSION="4.4.1"
GODOT_FILE="godot${GODOT_VERSION}"
GODOT_CMD="godotcmd"

GODOT_EXPORT_TEMPLATE="godot${GODOT_VERSION}_export_templte.tpz"
GODOT_TELEGRAM_WEB_EXPORT_PATH="$(pwd)/backend/telegram/static/game"
GODOT_PROJECT_LOCATION="$(pwd)/game/project.godot"
create_godot_cmd() {
	cd .cache
	if [[ -e "$GODOT_CMD" ]]; then
		echo "godot executable Exists"	
	else
		echo "No Godot executanble found"
		if ! [[ -e "${GODOT_FILE}.zip" ]]; then
			wget -O ${GODOT_FILE}.zip https://github.com/godotengine/godot-builds/releases/download/${GODOT_VERSION}-stable/Godot_v${GODOT_VERSION}-stable_linux.x86_64.zip
		else
			echo "not downloading ${GODOT_FILE}.zip already exist"
		fi

		unzip ${GODOT_FILE}.zip
		mv Godot_v${GODOT_VERSION}-stable_linux.x86_64  godotcmd
	fi
	cd ..
}

check_godot_export_template() {
	cd .cache
	if [[ -e "$GODOT_EXPORT_TEMPLATE" ]]; then
		echo "export template exists"	
	else
		echo "no export template exists"
		echo "expected file $(pwd)/${GODOT_EXPORT_TEMPLATE}"
		cd ..
		exit 1
	fi
	if ! [[ -e "$HOME/.local/share/godot/export_templates/${GODOT_VERSION}.stable/" ]] then
		unzip ${GODOT_EXPORT_TEMPLATE} -d ~/.local/share/godot/export_templates/${GODOT_VERSION}.stable
		mv ~/.local/share/godot/export_templates/${GODOT_VERSION}.stable/templates/* ~/.local/share/godot/export_templates/${GODOT_VERSION}.stable/
	else
		echo "template is unziped"
	fi
	cd ..
}

export_for_telegram_web() {
	echo "starting to export Web Game"
	mkdir -v -p $GODOT_TELEGRAM_WEB_EXPORT_PATH
	./.cache/godotcmd $GODOT_PROJECT_LOCATION --quiet --headless --export-release TelegramWeb "${GODOT_TELEGRAM_WEB_EXPORT_PATH}/mini.html"
	echo "Telegram Web Game Exported"
}

generate_sqlc() {
	echo "Generating..."
	cd ./backend/
	sqlc generate
	cd ..
	echo "OK"	
}

build_frontend() {
	echo "Building frontend"
	cd ./frontend/
	npm run build
	cd ..
}

COMMAND=$1

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
		generate_sqlc
		echo "Starting Hokm Bacbkend"
		cd ./backend/
		HOKM4_HTTP_PORT=4445 \
		  go run ./cmd/hokm4/.
		echo "The End"
		;;

	run_bot)
		generate_sqlc
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
	godot)
		create_godot_cmd
		check_godot_export_template
		build_frontend
		export_for_telegram_web
		;;

	*)
		echo "command '$COMMAND' is Unknown"
esac

exit 0
