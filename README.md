# tix

## Hokm Game

GG Stack:
    Godot
    Golang

## Deps:
- Sqlc: Sql Compiler
- Taskfile: Autmating some stuff
- Chi: routing only
- Gorrilla/websocket
- joho/godotenv: reading  .env files
- pascaldekloe/jwt

## Backend:
contains two service User & Hokm4
### User:
- profiles
- authentication & authorization
### Hokm4:
- websocket and game logic

# Run on your machine

### Requierments:
- install sqlc
- install task
- install docker

## Running:
```bash
task caddy
task run_user
task run_hokm_4
```
