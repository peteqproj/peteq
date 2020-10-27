.PHONY: build
build:
	go build -o ./dist/peteq ./cmd/server

.PHONY: dependency-update
dependency-update:
	go mod download

.PHONY: run
run:
	./dist/peteq

.PHONY: mock-all
mock-all:
	go get github.com/vektra/mockery/v2/.../
	mockery --all --inpackage

.PHONY: test
test:
	./hack/test.sh