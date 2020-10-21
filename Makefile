.PHONY: build
build:
	go build -o ./dist/peteq ./cmd/server

.PHONY: dependency-update
dependency-update:
	go mod download

.PHONY: run
run:
	./dist/peteq

.PHONY: run-docker
run-docker:
	env > .env
	docker run -p 8080:8080 -d --env-file ./.env peteqproj/peteq 
	rm .env