AUTOTESTS_BINARY=./bin/metricstest

.PHONY: get-autotests-arm64
get-autotests-arm64: export AUTOTESTS_TAG=`git ls-remote --tags --sort="-version:refname" https://github.com/Yandex-Practicum/go-autotests | head -n 1 | awk '{print $2}' | sed 's/refs\/tags\///' | grep -o 'v.*'`
get-autotests-arm64:
	@echo "downloading autotests binary runner"
	@echo "latest version: $(AUTOTESTS_TAG)"
	@wget https://github.com/Yandex-Practicum/go-autotests/releases/download/$(AUTOTESTS_TAG)/metricstest-darwin-arm64 -O ./bin/metricstest
	@chmod +x ./bin
	@echo "autotests binary runner downloaded"

.PHONY: run-autotests

ITER := `git rev-parse --abbrev-ref HEAD/iter`
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

run-autotests:
	@echo "running autotests with $(ARGS)"
	@if [ ! -f $(AUTOTESTS_BINARY) ]; then \
		echo "autotests binary not found, downloading..."; \
		$(MAKE) get-autotests-arm64; \
	fi
	@$(AUTOTESTS_BINARY) -test.v -test.run=^TestIteration$(ITER)

.PHONY: update-autotests
update-autotests:
	@echo "fetch autotests from template repo"
	@git fetch template && git checkout template/v2 .github
	@echo "autotests updated"

.PHONY: build-agent
build-agent:
	@echo "building agent"
	@cd ./cmd/agent && go build -o agent *.go

.PHONY: build-server
build-server:
	@echo "building server"
	@cd ./cmd/server && go build -o server *.go
