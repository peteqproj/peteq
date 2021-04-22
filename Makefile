.PHONY: build
build:
	make gen-openapi
	make build-dev-cli
	make gen-domain
	make gen-repo 
	make compile

.PHONY: compile
compile:
	make build-dev-cli
	go build -o ./dist/peteq-server cmd/server/main.go
	go build -o ./dist/peteq cmd/peteq-cli/main.go

.PHONY: build-dev-cli
build-dev-cli:
	go build -o ./dist/peteq-dev cmd/dev-cli/main.go

.PHONY: gen-domain
gen-domain:
	./dist/peteq-dev create aggregate --package task --schema manifests/task/task.json
	./dist/peteq-dev create aggregate --package user --schema manifests/user/user.json
	./dist/peteq-dev create aggregate --package list --schema manifests/list/list.json
	./dist/peteq-dev create aggregate --package project --schema manifests/project/project.json
	./dist/peteq-dev create aggregate --package sensor --schema manifests/sensor/sensor.json
	./dist/peteq-dev create aggregate --package automation --schema manifests/automation/automation.json --schema manifests/automation/sensor.binding.json

.PHONY: gen-repo
gen-repo:
	./dist/peteq-dev create repo --repo manifests/task/repo.yaml
	./dist/peteq-dev create repo --repo manifests/user/repo.yaml
	./dist/peteq-dev create repo --repo manifests/list/repo.yaml
	./dist/peteq-dev create repo --repo manifests/project/repo.yaml
	./dist/peteq-dev create repo --repo manifests/sensor/repo.yaml
	./dist/peteq-dev create repo --repo manifests/automation/repo.yaml
	gofmt -w -s .
	# go get golang.org/x/tools/cmd/goimports
	goimports -w .


.PHONY: dependency-update
dependency-update:
	go mod download


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

.PHONY: gen-openapi-client
gen-openapi-client:
	docker run --rm -v "${PWD}:/local" --user "$(id -u):$(id -g)" openapitools/openapi-generator-cli generate -i /local/docs/swagger.yaml -g go -o /local/pkg/client -p=isGoSubmodule=true -p=packageName=client
	sudo rm -rf pkg/client/api pkg/client/docs pkg/client/go.* pkg/client/git_push.sh pkg/client/README.md pkg/client/.openapi-generator pkg/client/.gitignore pkg/client/.openapi-generator-ignore pkg/client/.travis.yml
	gofmt -w -s .