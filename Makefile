docker_compose_bin := docker compose

build-all:
	cd cart && make build
	cd loms && make build
	cd notifications && make build

run-all: build-all
	$(docker_compose_bin) up --build

check:
	cd cart && make check
	cd loms && make check
	#cd notifications && make check
