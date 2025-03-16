.PHONY: test int-test

int-test: ## Run integration tests
	INT_TESTS=1 CONFIG=config.yaml go test ./tests/...

test: ## run unit tests
	go test ./...