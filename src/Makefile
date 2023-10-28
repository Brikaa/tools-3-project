docker_compose_build = docker-compose -f docker-compose.build.yaml
docker_compose_run = docker-compose -f docker-compose.run.yaml

build:
	$(docker_compose_build) up --build && $(docker_compose_build) rm -fsv

run:
	$(docker_compose_run) up --build -d

logs:
	$(docker_compose_run) logs -f

stop:
	$(docker_compose_run) down

dev:
	make stop && \
	make build && \
	make run && \
	make logs
