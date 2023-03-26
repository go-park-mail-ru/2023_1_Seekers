build:
	docker-compose up --build

prune:
	docker system prune

cov:
	./scripts/coverage.sh