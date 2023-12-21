SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.yml

# Bold
BCYAN=\033[1;36m
BBLUE=\033[1;34m

# No color (Reset)
NC=\033[0m

.DEFAULT_GOAL := help

.PHONY: s3-up
s3-up: ## Start S3 service
	${DOCKER_COMPOSE_PREFIX} up -d minio
	@echo "Waiting for Minio to become healthy..."
	@until docker-compose exec -T minio sh -c "curl -f http://localhost:9000/minio/health/live > /dev/null 2>&1"; do \
		sleep 1; \
		echo "Minio is not healthy yet, retrying..."; \
	done
	@printf "Minio is now healthy!\n\n"

.PHONY: s3-down
s3-down: ## Stop S3 service
	${DOCKER_COMPOSE_PREFIX} rm -fsv minio

.PHONY: clean	
clean: ## Cleanup
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

.PHONY: test
test: ## Run tests
ifndef GITHUB_ACTIONS
	$(MAKE) s3-up
endif
	export IS_LOCAL=false; \
	go vet ./...; \
	go test ./... -cover
ifndef GITHUB_ACTIONS
	$(MAKE) s3-down
endif

.PHONY: help
help: ## Disply this help
		@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BCYAN)%-18s$(NC)%s\n", $$1, $$2}'