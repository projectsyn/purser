MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

DOCKER_CMD   ?= docker
DOCKER_ARGS  ?= run --rm --user "$$(id -u)" -v "$${PWD}:/src" --workdir /src

YAML_FILES      ?= $(shell find . -not -path './vendor/*' -type f \( -name '*.yaml' -or -name '*.yml' \) )
YAMLLINT_ARGS   ?= --no-warnings
YAMLLINT_CONFIG ?= .yamllint.yml
YAMLLINT_IMAGE  ?= docker.io/cytopia/yamllint:latest
YAMLLINT_DOCKER ?= $(DOCKER_CMD) $(DOCKER_ARGS) $(YAMLLINT_IMAGE)

VALE_CMD  ?= $(DOCKER_CMD) $(DOCKER_ARGS) --volume "$${PWD}"/docs/modules:/pages vshn/vale:2.1.1
VALE_ARGS ?= --minAlertLevel=error --config=/pages/ROOT/pages/.vale.ini /pages


.PHONY: all
all: lint

.PHONY: lint
lint: lint_yaml lint_adoc

.PHONY: lint_yaml
lint_yaml: $(YAML_FILES)
	$(YAMLLINT_DOCKER) -f parsable -c $(YAMLLINT_CONFIG) $(YAMLLINT_ARGS) -- $?

.PHONY: lint_adoc
lint_adoc:
	$(VALE_CMD) $(VALE_ARGS)
