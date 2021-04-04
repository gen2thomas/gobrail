.PHONY: run test fmt

excluding_vendor := $(shell go list ./... | grep -v /vendor/)

# Run latest working level
run:
	go run cmd/main_board.go

# Run tests on all non-vendor directories
test:
	go test -v $(excluding_vendor)

# Correct format errors
fmt:
	go fmt ./...

# goplantuml must be installed first with "GO111MODULE=off" prefix
plantuml:
	goplantuml -recursive -hide-fields -hide-methods ./ > docs/new.puml