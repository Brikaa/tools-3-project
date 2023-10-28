docker_compose_build = docker-compose -f docker-compose.build.yaml
docker_compose_run = docker-compose -f docker-compose.run.yaml
define with_env =
	export $(shell cat .env) && \
	$1
endef

dev:
	make stop
	make build
	make run
	make logs

build:
	$(docker_compose_build) up --build && $(docker_compose_build) rm -fsv

run:
	./create-env.sh
	$(docker_compose_run) up --build -d

logs:
	$(docker_compose_run) logs -f

stop:
	$(docker_compose_run) down

migrate:
	$(call with_env,\
	$(docker_compose_run) exec -it database /bin/sh -c "mysql app -uroot -p${MYSQL_ROOT_PASSWORD} < /initial.sql")

mysql:
	$(call with_env,\
	$(docker_compose_run) exec -it database mysql -hdatabase -u${MYSQL_USER} -p${MYSQL_PASSWORD} -Dapp)
