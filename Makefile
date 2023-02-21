# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# compose-up: ### Run docker-compose
# 	docker-compose up --build -d postgres rabbitmq && docker-compose logs -f
# .PHONY: compose-up

# compose-up-integration-test: ### Run docker-compose with integration test
# 	docker-compose up --build --abort-on-container-exit --exit-code-from integration
# .PHONY: compose-up-integration-test

# compose-down: ### Down docker-compose
# 	docker-compose down --remove-orphans
# .PHONY: compose-down

swag: ### swag init
	swag init -g handlers/httpserver/router.go
.PHONY: swag


run: fmt swag ### regenerate swag docs, check module and run go code
	go mod tidy
	go mod download
	go run .
.PHONY: run

fmt: ### format swag docs, go mod and code
	go mod tidy 
	go fmt .
.PHONY: fmt

lint: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test: ### run test
	go test -v github.com/vano2903/service-template/controller/tests
.PHONY: test

update: ### update dependencies
	go get -u
	go mod tidy
.PHONY: update

docker: ### build and run docker image
	docker build -t service-template .
	docker run -p 8080:8080 service-template
.PHONY: docker

prep: fmt swag lint test ### run all checks before commit
.PHONY: prep