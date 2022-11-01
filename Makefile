.PHONY: all run build buildarm5 buildarm7 test style

mainfile := cmd/main_daemon.go

all: build buildarm5 buildarm7

# Run latest working level
run:
	go run $(mainfile)

build: 
	go build -o ./output/gobrail $(mainfile)

buildarm5:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ./output/gobrail_raspi $(mainfile)

buildarm7:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./output/gobrail_tinker $(mainfile)

# Run tests on all non-vendor directories
test:
	go test -v ./... -cover  -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Correct format errors amd check linting
style:
	go fmt ./...
	golint ./...

# goplantuml must be installed first with "GO111MODULE=off" prefix
plantuml:
	goplantuml -recursive -hide-fields -hide-methods ./ > docs/new.puml
