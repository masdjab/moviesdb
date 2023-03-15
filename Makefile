DOCKER_COMMAND=docker-compose -f ./docker/docker-compose.yml --env-file docker/.env

copy-config:
	cp config.yml.sample config.yml

docker.start:
	${DOCKER_COMMAND} up -d

docker.stop:
	${DOCKER_COMMAND} down

migrate:
	go run main.go migrate

test:
	go test ./... -covermode=atomic -coverprofile=coverage.out
	go tool cover -func coverage.out

start:
	go run main.go
