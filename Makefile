# note: call scripts from /deploy

# project name
PROJECTNAME=$(shell basename "$(PWD)")
BUILD_VERSION	:= 1.0.0
BUILD_TIME		:= $(shell date "+%F %T")
BUILD_NAME		:= supervisord
SOURCE			:= ./
TARGET_DIR		:= /${BUILD_NAME}

# project path
ROOT=$(shell pwd)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

tidy:
	@echo "use mod tidy"
	@export GO111MODULE=on
	@export GOPROXY=https://goproxy.io
	@go mod tidy

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_NAME} ${SOURCE}
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o push bin/push.go
mac:
	go build -o push bin/push.go