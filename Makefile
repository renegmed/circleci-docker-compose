up:
	docker-compose up --build -d 
.PHONY: up

test:
	$(MAKE) -C integrationtest test-originator
	$(MAKE) -C integrationtest test-salesorders
	$(MAKE) -C integrationtest test-joborders
