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

build-binaries:
	go build cmd/api/main.go
	go build cmd/auth/main.go
	go build cmd/file_storage/main.go
	go build cmd/auth/main.go
	go build cmd/mail/main.go
	go build cmd/user/main.go

build-file_storage:
	docker-compose -f docker-compose-prod.yml up -d --build --remove-orphans file_storage

build-user:
	docker-compose -f docker-compose-prod.yml up -d --remove-orphans user

build-auth:
	docker-compose -f docker-compose-prod.yml up -d --remove-orphans auth

build-mail:
	docker-compose -f docker-compose-prod.yml up -d --remove-orphans mail

build-api:
	docker-compose -f docker-compose-prod.yml up -d --remove-orphans api

build-prod-cd:
	sudo systemctl stop nginx.service || true
	sudo systemctl stop postgresql || true
	mkdir -p -m 777 logs/app
	mkdir -p -m 777 logs/postgres
	docker-compose -f docker-compose-prod.yml up -d --remove-orphans
	sudo cp ./nginx/nginx.prod.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

build-prod:
	sudo systemctl stop nginx.service || true
	sudo systemctl stop postgresql || true
	mkdir -p -m 777 logs/app
	mkdir -p -m 777 logs/postgres
	docker-compose -f docker-compose-prod.yml up -d --build --remove-orphans
	sudo cp ./nginx/nginx.prod.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

build-dev-env:
	sudo systemctl stop nginx.service || sudo launchctl nginx || true
	sudo systemctl stop postgresql || sudo launchctl stop postgresql-15 ||  true
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app
	make docker-stop-back
	make docker-rm-volumes || true
	docker-compose up -d --build --remove-orphans db
	sudo docker exec -it postgres psql -U ${POSTGRES_USER} -W ${POSTGRES_PASSWORD} -d ${POSTGRES_DB} -a -f /setup_user_sql/000_setup_users.sql -v db_mail_service_user_pw="'${DB_MAIL_SERVICE_USER_PW}'" -v db_user_service_user_pw="'${DB_USER_SERVICE_USER_PW}'"
	docker-compose up -d --remove-orphans
#	sudo cp ./nginx/nginx.conf /etc/nginx/nginx.conf
#	sudo systemctl restart nginx

run-dev:
	@make build-dev-env
	@make run-all-services POSTGRES_HOST=127.0.0.1 REDIS_HOST=127.0.0.1

run-all-services:
	make -j 5 run-api-service run-auth-service run-file_storage-service run-mail-service run-user-service

run-auth-service:
	go run ./cmd/auth/main.go -config=./cmd/config/debug.yml

run-api-service:
	go run ./cmd/api/main.go -config=./cmd/config/debug.yml

run-user-service:
	go run ./cmd/user/main.go -config=./cmd/config/debug.yml

run-mail-service:
	go run ./cmd/mail/main.go -config=./cmd/config/debug.yml

run-file_storage-service:
	go run ./cmd/file_storage/main.go -config=./cmd/config/debug.yml

run-smtp-server:
	sudo go run cmd/smtp/main.go -config=./cmd/config/debug.yml

docker-prune:
	@bash -c 'docker system prune -f'

docker-stop-back:
	docker container stop grafana || true
	docker container stop prometheus || true
	docker container stop 2023_1_seekers_api_1 || true
	docker container stop 2023_1_seekers_mail_1 || true
	docker container stop 2023_1_seekers_auth_1 || true
	docker container stop 2023_1_seekers_user_1 || true
	docker container stop admin_db || true
	docker container stop postgres || true
	docker container stop cache || true
	docker container stop redis || true
	docker container stop node_exporter || true
	make docker-prune

docker-stop-all:
	@bash -c "docker kill $(shell eval docker ps -q)"
	@bash -c "docker rm $(shell eval docker ps -a -q)"
	@make docker-prune

docker-rm-volumes:
	@bash -c "docker volume rm $(shell eval docker volume ls -qf dangling=true)"

docker-rm-images:
	@bash -c "docker rmi $(shell eval docker images -q)"

cov:
	mkdir -p ${COV_DIR}
	go test -race -coverpkg=${INTERNAL_PKG} -coverprofile ${COV_DIR}/cover.out ${INTERNAL_PKG}; cat ${COV_DIR}/cover.out | fgrep -v "test.go" | fgrep -v "register.go"| fgrep -v "docs" | fgrep -v ".pb.go" | fgrep -v "mock" | fgrep -v "config" | fgrep -v "_easyjson.go" > ${COV_DIR}/cover2.out
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

perm-dirs:
	sudo chmod -R 777 ./

lint:
	golangci-lint run -c ./.golangci.yml

test:
	go test ./...
