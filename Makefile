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

.PHONY: run-autotests

ITER := `git rev-parse --abbrev-ref HEAD | grep -o '[0-9]\+'`

run-autotests:
	@echo "running autotests for iteration $(ITER)"
	@if [ ! -f $(AUTOTESTS_BINARY) ]; then \
		echo "autotests binary not found, downloading..."; \
		$(MAKE) get-autotests-arm64; \
	fi
	@$(AUTOTESTS_BINARY) -test.v -test.run=^TestIteration$(ITER)$ -binary-path=./cmd/server/server -agent-binary-path=./cmd/server/agent

.PHONY: update-autotests
update-autotests:
	@echo "fetch autotests from template repo"
	@git fetch template && git checkout template/v2 .github
	@echo "autotests updated"

.PHONY: build-agent
build-agent:
	@echo "building agent"
	@go build -o cmd/server/agent cmd/agent/main.go

.PHONY: build-server
build-server:
	@echo "building server"
	@go build -o cmd/server/server cmd/server/main.go
