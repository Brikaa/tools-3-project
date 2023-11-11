-include .env
docker_compose_test = docker compose -f docker-compose.test.yaml

dev:
	./create-env.sh
	make cache
	make run
	make logs

cache:
	docker build --target build -t backend_builder ./backend
	docker build --target build -t frontend_builder ./frontend

run:
	docker compose up --build -d

logs:
	docker compose logs -f

stop:
	docker compose down

migrate:
	docker compose exec -it database /bin/sh -c "mysql app -uroot -p${MYSQL_ROOT_PASSWORD} < /initial.sql"

mysql:
	docker compose exec -it database mysql -hdatabase -u${MYSQL_USER} -p${MYSQL_PASSWORD} -Dapp

test:
	make migrate
	$(docker_compose_test) up --build && $(docker_compose_test) rm -fsv
