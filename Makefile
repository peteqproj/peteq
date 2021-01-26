.PHONY: build
build:
	go build -o ./dist/peteq main.go

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
	go get github.com/vektra/mockery/v2/.../
	mockery --all --inpackage

.PHONY: test
test:
	./hack/test.sh