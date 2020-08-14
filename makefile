.PHONY: cli
cli:
	docker-compose run --rm --service-ports app sh

.PHONY: compile
compile:
	CGO_ENABLED=0 \
	GO111MODULES=on \
	go build \
		-a \
		-o ./bin/server \
		-ldflags '-extldflags -static' \
		cmd/server.go
