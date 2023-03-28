build:
	docker-compose up --build

no-env-build:
	docker-compose -f ./docker-compose-NO_ENV.yml up --build

prune:
	docker system prune

cov:
	./scripts/coverage.sh