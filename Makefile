build:
	docker-compose up --build

no-env-build:
	docker-compose -f ./docker-compose-NO_ENV.yml up --build

docker-prune:
	docker system prune

docker-prune-volumes:
	docker volume rm $(docker volume ls -qf dangling=true)

cov:
	./scripts/coverage.sh
	
clean_logs:
	sudo rm -rf logs/*
	mkdir -m 777 logs/postgres
	mkdir -m 777 logs/app

