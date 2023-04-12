INTERNAL_PKG = ./internal/...
ALL_PKG = ./internal/... ./pkg/...
COV_DIR = scripts/result_cover

build:
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app
	docker-compose up -d --build
	sudo cp ./nginx/nginx.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

build-prod:
	mkdir -p -m 777 logs/app
	mkdir -p -m 777 logs/postgres
	docker-compose -f docker-compose-prod.yml up -d --build
	sudo cp ./nginx/nginx.prod.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

docker-prune:
	docker system prune

docker-prune-volumes:
	docker volume rm $(docker volume ls -qf dangling=true)

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
	swag init -g cmd/main.go -o docs

generate:
	go generate ${ALL_PKG}