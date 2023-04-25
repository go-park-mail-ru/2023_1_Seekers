INTERNAL_PKG = ./internal/...
ALL_PKG = ./internal/... ./pkg/...
COV_DIR = scripts/result_cover

include .env
export

build:
	sudo systemctl stop nginx.service || true  
	sudo systemctl stop postgresql || true 
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app
	docker-compose up -d --build --remove-orphans

build-prod:
	sudo systemctl stop nginx.service
	sudo systemctl stop postgresql
	mkdir -p -m 777 logs/app
	mkdir -p -m 777 logs/postgres
	docker-compose -f docker-compose-prod.yml up -d --build --remove-orphans
	sudo cp ./nginx/nginx.prod.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

build-dev-env:
	sudo systemctl stop nginx.service || true 
	sudo systemctl stop postgresql || true 
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app
	docker-compose -f docker-compose-dev.yml up -d --build --remove-orphans

run-dev:
	@make build-dev-env
	@make run-all-services POSTGRES_HOST=127.0.0.1 REDIS_HOST=127.0.0.1

run-all-services:
	make -j 5 run-api-service run-auth-service run-file_storage-service run-mail-service run-user-service

run-auth-service:
	go run ./cmd/auth/main.go -config=./cmd/config/dev.yml

run-api-service:
	go run ./cmd/api/main.go -config=./cmd/config/dev.yml

run-user-service:
	go run ./cmd/user/main.go -config=./cmd/config/dev.yml

run-mail-service:
	go run ./cmd/mail/main.go -config=./cmd/config/dev.yml

run-file_storage-service:
	go run ./cmd/file_storage/main.go -config=./cmd/config/dev.yml

docker-prune:
	@bash -c 'docker system prune'

docker-stop:
	@bash -c "docker kill $(shell eval docker ps -q)"
	@bash -c "docker rm $(shell eval docker ps -a -q)"
	@make docker-prune

docker-rm-volumes:
	@bash -c "docker volume rm $(shell eval docker volume ls -qf dangling=true)"

docker-rm-images:
	@bash -c "docker rmi $(shell eval docker images -q)"

cov:
	mkdir -p ${COV_DIR}
	go test -race -coverpkg=${INTERNAL_PKG} -coverprofile ${COV_DIR}/cover.out ${INTERNAL_PKG}; cat ${COV_DIR}/cover.out | fgrep -v "test.go" | fgrep -v "register.go"| fgrep -v "docs" | fgrep -v ".pb.go" | fgrep -v "mock" > ${COV_DIR}/cover2.out
	go tool cover -func ${COV_DIR}/cover2.out
	go tool cover -html ${COV_DIR}/cover2.out -o ${COV_DIR}/coverage.html

clean_logs:
	sudo rm -rf logs/*
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app

doc:
	swag init -g cmd/api/main.go -o docs

generate:
	go generate ${ALL_PKG}