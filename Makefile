up:
	docker-compose up --build -d 
.PHONY: up

test:
	SERVER=7070 NATS_SERVER_ADDR=nats://demo.nats.io CLIENT_ID=tester go test integrationtest/*.go 
