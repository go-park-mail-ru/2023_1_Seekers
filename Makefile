build:
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app
	docker-compose up -d
	sudo cp ./nginx/nginx.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

build-prod:
	mkdir -p -m 777 logs/app
	mkdir -p -m 777 logs/postgres
	docker-compose -f docker-compose-prod.yml up -d
	sudo cp ./nginx/nginx.prod.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

docker-prune:
	docker system prune

docker-prune-volumes:
	docker volume rm $(docker volume ls -qf dangling=true)

cov:
	mkdir -p scripts/result_cover
	go test -race -coverpkg=./... -coverprofile scripts/result_cover/cover.out ./...; cat scripts/result_cover/cover.out | fgrep -v "test.go" | fgrep -v "register.go"| fgrep -v "docs" | fgrep -v ".pb.go" | fgrep -v "mock" > scripts/result_cover/cover2.out
	go tool cover -func scripts/result_cover/cover2.out
	go tool cover -html scripts/result_cover/cover2.out -o scripts/result_cover/coverage.html

clean_logs:
	sudo rm -rf logs/*
	mkdir -p -m 777 logs/postgres
	mkdir -p -m 777 logs/app

