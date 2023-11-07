docker_compose_bin := docker compose

run-all:
	$(docker_compose_bin) up --build

check:
	cd cart && make check
	cd loms && make check
	cd notifications && make check

up-kafka:
	$(docker_compose_bin) -f docker-compose-kafka.yml up

migrate:
	$(docker_compose_bin) exec cart goose -dir internal/app/storage/migrations up
	$(docker_compose_bin) exec loms goose -dir internal/app/storage/migrations up