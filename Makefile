FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTHENTICATION_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_authentication build_logger build_mail build_listener build_front
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	cd ./broker-service && docker build -f broker-service.dockerfile -t petarnenov/broker-service .
	# docker push petarnenov/broker-service
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} ./cmd/api
	cd ./listener-service && docker build -f listener-service.dockerfile -t petarnenov/listener-service .
	docker push petarnenov/listener-service
	@echo "Done!"

## build_authentication: builds the authentication binary as a linux executable
build_authentication:
	@echo "Building authentication binary..."
	cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHENTICATION_BINARY} ./cmd/api
	cd ./authentication-service && docker build -f authentication-service.dockerfile -t petarnenov/authentication-service .
	docker push petarnenov/authentication-service
	@echo "Done!"

## build_authentication: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	cd ./logger-service && docker build -f logger-service.dockerfile -t petarnenov/logger-service .
	docker push petarnenov/logger-service
	@echo "Done!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	cd ./mail-service && docker build -f mail-service.dockerfile -t petarnenov/mail-service .
	docker push petarnenov/mail-service
	@echo "Done!"

## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd ./front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	cd ./front-end && docker build -f front-end-service.dockerfile -t petarnenov/front-end-service .
	docker push petarnenov/front-end-service
	@echo "Done!"

## build caddy
build_caddy:
	@echo "Building caddy..."
	docker build -f caddy.dockerfile -t petarnenov/caddy .
	docker push petarnenov/caddy
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ./front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## restart: restarts all
restart: stop down up_build
	@echo "Restarted all!"

## create docker swarm
swarm_up:
	docker swarm init

swarm_down:
	docker swarm leave --force

create_proto:
	cd logger-service/logs
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto

deploy_stack:
	docker stack deploy --compose-file docker-compose.yaml microservices

rm_stack:
	docker stack rm microservices

## docker service update --image petarnenov/front-end-service microservices_front-end-service

