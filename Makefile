#export_env:
#	./scripts/export_env.sh
include .env

build:
	./scripts/export_env.sh && docker-compose up -d --build
prune:
	docker system prune