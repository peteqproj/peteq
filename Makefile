.PHONY: build
build:
	make gen-openapi && go build -o ./dist/peteq main.go

@.PHONY: build-cron-watcher
build-cron-watcher:
	go build -o ./dist/peteq-cron-wacher ./cmd/cron

.PHONY: dependency-update
dependency-update:
	go mod download

.PHONY: run
run:
	./dist/peteq

.PHONY: run-cron-watcher
run-cron-watcher:
	PORT=8082 ./dist/peteq-cron-wacher

.PHONY: mock-all
mock-all:
	docker pull vektra/mockery:latest
	docker run --workdir=/app -v $(PWD):/app vektra/mockery:latest --all --inpackage	

.PHONY: gen-openapi
gen-openapi:
	docker run --workdir=/app -v $(PWD):/app peteqproj/openapi swag init -g pkg/server/openapi.go

.PHONY: test
test:
	./hack/test.sh