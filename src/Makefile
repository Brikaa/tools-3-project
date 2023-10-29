-include .env
docker_compose_build = docker-compose -f docker-compose.build.yaml
docker_compose_run = docker-compose -f docker-compose.run.yaml

dev:
	./create-env.sh
	make stop
	make build
	make run
	make logs

build:
	$(docker_compose_build) up --build && $(docker_compose_build) rm -fsv

run:
	$(docker_compose_run) up --build -d

logs:
	$(docker_compose_run) logs -f

stop:
	$(docker_compose_run) down

migrate:
	$(docker_compose_run) exec -it database /bin/sh -c "mysql app -uroot -p${MYSQL_ROOT_PASSWORD} < /initial.sql"

mysql:
	$(docker_compose_run) exec -it database mysql -hdatabase -u${MYSQL_USER} -p${MYSQL_PASSWORD} -Dapp
