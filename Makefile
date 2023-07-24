GO := go
PACKAGES=$(shell go list ./...)

# run the service in docker
start:
	docker-compose up -d
.PHONY: run

stop:
	docker-compose down
.PHONY: run

ask:
	docker-compose run client
.PHONY: ask

run: start ask stop

test:
	$(GO) test -p 1 $(PACKAGES)
.PHONY: test
