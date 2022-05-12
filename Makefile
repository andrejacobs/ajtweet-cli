# TODO: Add copyright notice
# TODO: Describe what make options you have

EXECUTABLE_BASE_NAME=ajtweet
BUILD_OUTPUT_DIR=build

CURRENT_OS=$(shell uname -s)
CURRENT_CPU_ARCH=$(shell uname -p)

.PHONY: all
all: clean info build test

# Fetch dependencies
.PHONY: deps
deps:
	go get -u github.com/spf13/cobra@latest

# Run unit-tests
.PHONY: test
test:
	go test -v ./...

# Build executable for the current platform
.PHONY: build
build:
	$(eval CURRENT_OUTPUT_DIR := ${BUILD_OUTPUT_DIR}/${CURRENT_OS}-${CURRENT_CPU_ARCH})
	$(eval CURRENT_EXECUTABLE := ${CURRENT_OUTPUT_DIR}/${EXECUTABLE_BASE_NAME})
	@echo 
	mkdir -p ${CURRENT_OUTPUT_DIR}
	go build -o ${CURRENT_EXECUTABLE} main.go

# Remove build output
.PHONY: clean
clean:
	go clean
	rm -rf ${BUILD_OUTPUT_DIR}

# Display info about the current platform and configurations etc.
.PHONY: info
info:
	@echo "OS: ${CURRENT_OS}"
	@echo "CPU architecture: ${CURRENT_CPU_ARCH}"
