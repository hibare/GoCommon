SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.yml

all: s3-up

s3-up:
	${DOCKER_COMPOSE_PREFIX} up -d minio

s3-down:
	${DOCKER_COMPOSE_PREFIX} rm -fsv minio
	
clean: 
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

test: 
ifndef GITHUB_ACTIONS
	$(MAKE) s3-up
endif
	go test ./... -cover
ifndef GITHUB_ACTIONS
	$(MAKE) s3-down
endif

.PHONY = all clean test s3-up s3-down