AUTOTESTS_BINARY=./bin/metricstest

.PHONY: get-autotests-arm64
get-autotests-arm64: export AUTOTESTS_TAG=`git ls-remote --tags --sort="-version:refname" https://github.com/Yandex-Practicum/go-autotests | head -n 1 | awk '{print $2}' | sed 's/refs\/tags\///' | grep -o 'v.*'`
get-autotests-arm64:
	@echo "downloading autotests binary runner"
	@echo "latest version: $(AUTOTESTS_TAG)"
	@mkdir -p ./bin
	@wget https://github.com/Yandex-Practicum/go-autotests/releases/download/$(AUTOTESTS_TAG)/metricstest-darwin-arm64 -O ./bin/metricstest
	@chmod +x ./bin/metricstest
	@chmod 755 ./bin/metricstest
	@echo "autotests binary runner downloaded"

.PHONY: update-autotests
update-autotests:
	@echo "fetch autotests from template repo"
	@git fetch template && git checkout template/v2 .github
	@echo "autotests updated"

.PHONY: build-agent
build-agent:
	@echo "building agent"
	@go build -o cmd/agent/agent cmd/agent/main.go

.PHONY: build-server
build-server:
	@echo "building server"
	@go build -o cmd/server/server cmd/server/main.go

.PHONY: run-agent
run-agent:
	@echo "running agent"
	@go run cmd/agent/main.go

.PHONY: run-server
run-server:
	@echo "running server"
	@go run cmd/server/main.go

.PHONY: test
test:
	@echo "running tests"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out | grep total
	@rm -f coverage.out

.PHONY: run-autotests

ITER := `git rev-parse --abbrev-ref HEAD | grep -o '[0-9]\+'`
ITER_BRANCH := ^TestIteration$(ITER)$

include .env

run-autotests: build-server build-agent
	@echo "running autotests for iteration $(ITER)"
	@if [ ! -f $(AUTOTESTS_BINARY) ]; then \
		echo "autotests binary not found, downloading..."; \
		$(MAKE) get-autotests-arm64; \
	fi
	@$(AUTOTESTS_BINARY) -test.v -test.run=$(ITER_BRANCH) -binary-path=./cmd/server/server -agent-binary-path=./cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=metrics.json -database-dsn="host=localhost port=5432 user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable"

.PHONY: gen
gen:
	@go generate ./...