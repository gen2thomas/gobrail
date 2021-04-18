.PHONY: run build buildraspi test style

excluding_vendor := $(shell go list ./... | grep -v /vendor/)
mainfile := cmd/main_daemon.go

# Run latest working level
run:
	go run $(mainfile)

build: 
	go build -o ./output/gobrail $(mainfile)

buildraspi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ./output/gobrail_raspi $(mainfile)

# Run tests on all non-vendor directories
test:
	go test -v $(excluding_vendor) -cover  -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Correct format errors amd check linting
style:
	go fmt ./...
	golint ./...

# goplantuml must be installed first with "GO111MODULE=off" prefix
plantuml:
	goplantuml -recursive -hide-fields -hide-methods ./ > docs/new.puml