-include .env
current_dir := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

# Tools.
export TOOLS=$(current_dir)tools
export TOOLS_BIN=$(TOOLS)/bin
export PATH := $(TOOLS_BIN):$(PATH)
export GOBIN=$(TOOLS_BIN)

help:  ## - these instruction
	@echo "\nTARGETS:\n"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/\n\t/'
	@echo ""

.PHONY:
run: ## - run app
	@bash -c "trap 'exit 0' INT; bash run.sh" 

.PHONY:
test: ## - run tests
	go test -v ./...


.PHONY:
test-coverage: ## - run coverage
	bash scripts/run_coverage.sh


.PHONY:  
hurl_test: ## - run tests
	hurl --test --verbose --variable host=localhost:8089 scripts/requests/event_test.hurl

.PHONY:
install-tools:	## - install tools goose, golangci-lint, gofumpt
	bash scripts/install-tools.sh


.PHONY:
lint: ## - run linter
	$(TOOLS_BIN)/golangci-lint run

.PHONY:
build: ## - build app container
	docker build . -t targetads

.PHONY:
compose-up: ## - run docker compose
	docker-compose --env-file .env up -d

%:
	@: