MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

# all targets are phony
.PHONY: $(shell egrep -o ^[a-zA-Z_-]+: $(MAKEFILE_LIST) | sed 's/://')

help: ## Print this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

init: ## create database in order to execute examples
	mysql -uroot -e "CREATE DATABASE IF NOT EXISTS torm"
	mysql -uroot -e "CREATE TABLE IF NOT EXISTS torm.users (id INT PRIMARY KEY AUTO_INCREMENT, org_id INT, name VARCHAR(255), email VARCHAR(255), age INT)"
	mysql -uroot -e "CREATE TABLE IF NOT EXISTS torm.orgs (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), address VARCHAR(255), established_at DATETIME)"

clean: ## destruct database
	mysql -uroot -e "DROP DATABASE torm"
