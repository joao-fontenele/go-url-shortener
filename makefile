.PHONY: cli
cli:
	docker-compose exec app sh

.PHONY: start
start:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: logs
logs:
	docker logs -f --since 1h --tail 300 go-url-shortener

.PHONY: compile
compile:
	CGO_ENABLED=0 \
	GO111MODULES=on \
	go build \
		-a \
		-o ./bin/server \
		-ldflags '-extldflags -static' \
		cmd/server.go
