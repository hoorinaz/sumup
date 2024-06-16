run-all-services:
	docker-compose up -d

clean:
	docker-compose down

test:
	go test ./asset/... -v

run:
	docker-compose up

build:
	docker-compose build