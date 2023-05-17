BUILD_DIR ?= build/bin
MODULE = $(shell go list -m)
define LD_FLAGS
"-X '$(MODULE)/pkg/version.Commit=$$(git rev-parse HEAD)'\
 -X '$(MODULE)/pkg/version.Branch=$$(git rev-parse --abbrev-ref HEAD)'\
 -X '$(MODULE)/pkg/version.Tag=$$(git name-rev --tags --name-only HEAD)'"
endef

.PHONY: tools
tools:
	@go install github.com/vektra/mockery/v2@v2.26.1
	@go install github.com/goreleaser/goreleaser@v1.18.2

.PHONY: generate
generate:
	@go generate ./...

.PHONY: docs
docs:
	@npx @redocly/cli build-docs ./docs/api.yaml -o ./docs/docs.html

.PHONY: deps
deps:
	@go mod tidy

tests: test.clean test test.build

test:
	@echo "Run tests"
	@go test ./... -timeout 5m

test.clean:
	@echo "Clean test cache"
	@go clean -testcache

test.race:
	@echo "Run tests with race"
	@go test ./... -race -timeout 5m

test.build:
	@echo "Run build test"
	@go build -o /dev/null $(MODULE)/cmd/server

.PHONY: lint
lint:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:latest golangci-lint run -v --timeout 5m

.PHONY: precommit
precommit: generate deps lint

.PHONY: build run
build:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)
	@go build -ldflags=${LD_FLAGS} -v -o $(BUILD_DIR)/apiserver $(MODULE)/cmd/server

run: build
	./$(BUILD_DIR)/apiserver --config ./config/local.yml

# FIXME: compose.%.% i.e. compose.${file}.${command}
compose.infra.%:
	$(eval CMD = ${subst compose.infra.,,$(@)})
	scripts/compose.sh infra $(CMD)

compose.local.%:
	$(eval CMD = ${subst compose.local.,,$(@)})
	scripts/compose.sh local $(CMD)

.PHONY: dist.clean dist release.local release.dry release

release.local: deps
	export LD_FLAGS=$(LD_FLAGS) && goreleaser release --snapshot --clean

release.dry: deps
	export LD_FLAGS=$(LD_FLAGS) && goreleaser build --clean

release: deps
	export LD_FLAGS=$(LD_FLAGS) && goreleaser release --clean

.PHONY: relocate
relocate:
	@test ${TARGET} || ( echo ">> TARGET is not set. Use: make relocate TARGET=<target>"; exit 1 )
	@find . -type f -name '*.go' -exec sed -i '' -e 's/$(MODULE)/{NEW_MODULE}/g' {} \;
	# Complete...
	# NOTE: This does not update the git config nor will it update any imports of the root directory of this project.