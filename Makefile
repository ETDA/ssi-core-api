start:
	docker-compose up -d

start-build:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	make stop && make start


restart-app:
	docker-compose restart abci did vc_status


restart-build:
	make stop && make start-build

logs-abci:
	 docker-compose logs -f abci

logs-did:
	 docker-compose logs -f did

logs-vc_status:
	 docker-compose logs -f vc_status


exec-abci:
	docker exec -it abci bash

test:
	go test ./...

download-module:
	go mod download

install:
	export GOPRIVATE=gitlab.finema.co/finema/* && git config --global url."git@gitlab.finema.co:".insteadOf "https://gitlab.finema.co/" && go get
