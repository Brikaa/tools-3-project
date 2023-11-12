-include .env
docker_or_podman = . ./container && docker
docker_compose_test = $(docker_or_podman) compose -f docker-compose.test.yaml

dev:
	./create-env.sh
	make cache
	make run
	make logs

cache:
	$(docker_or_podman) build --target build -t t3p-backend-builder ./backend
	$(docker_or_podman) build --target build -t t3p-frontend-builder ./frontend

run:
	$(docker_or_podman) compose up --build -d

logs:
	$(docker_or_podman) compose logs -f

stop:
	$(docker_or_podman) compose down

migrate:
	$(docker_or_podman) compose exec -it t3p-database \
	/bin/sh -c "mysql app -uroot -p${MYSQL_ROOT_PASSWORD} < /initial.sql"

mysql:
	$(docker_or_podman) compose exec -it t3p-database mysql -h ${MYSQL_HOST} -u${MYSQL_USER} -p${MYSQL_PASSWORD} -Dapp

test:
	make migrate
	$(docker_compose_test) up --build && $(docker_compose_test) rm -fsv
