# Copyright © 2022 André Jacobs

EXECUTABLE_BASE_NAME=ajtweet
BUILD_OUTPUT_DIR=build
MODULE_NAME="github.com/andrejacobs/ajtweet-cli"

.PHONY: all
all: clean info build test

# Fetch dependencies
.PHONY: deps
deps:
	@echo "Fetching dependencies"
	go mod tidy -v

# Run unit-tests
.PHONY: test
test:
	@echo "Running unit-tests"
	go test -v ./...

# Build executable for the current platform
.PHONY: build
build: versioninfo
	@echo "Building for the current platform"
	$(eval CURRENT_OUTPUT_DIR := ${BUILD_OUTPUT_DIR}/current)
	$(eval CURRENT_EXECUTABLE := ${CURRENT_OUTPUT_DIR}/${EXECUTABLE_BASE_NAME})
	mkdir -p ${CURRENT_OUTPUT_DIR}
	go build -ldflags "${GO_LDFLAGS}" -o ${CURRENT_EXECUTABLE} main.go

# Build executables for all the supported platforms
.PHONY: buildall
buildall: build-darwin build-linux

.PHONY: build-darwin
build-darwin: versioninfo
	@echo "Building for Darwin"
	GOOS=darwin GOARCH=arm64 go build -ldflags "${GO_LDFLAGS}" -o ${BUILD_OUTPUT_DIR}/Darwin/arm64/${EXECUTABLE_BASE_NAME} main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "${GO_LDFLAGS}" -o ${BUILD_OUTPUT_DIR}/Darwin/amd64/${EXECUTABLE_BASE_NAME} main.go

.PHONY: build-linux
build-linux: versioninfo
	@echo "Building for Linux"
	GOOS=linux GOARCH=arm64 go build -ldflags "${GO_LDFLAGS}" -o ${BUILD_OUTPUT_DIR}/Linux/arm64/${EXECUTABLE_BASE_NAME} main.go
	GOOS=linux GOARCH=amd64 go build -ldflags "${GO_LDFLAGS}" -o ${BUILD_OUTPUT_DIR}/Linux/amd64/${EXECUTABLE_BASE_NAME} main.go

# Build and run
.PHONY: run
run: build
	./${CURRENT_EXECUTABLE}

# Remove build output
.PHONY: clean
clean:
	@echo "Cleaning build output"
	go clean
	rm -rf ${BUILD_OUTPUT_DIR}

# Gather info about the current platform and version for the app
.PHONY: versioninfo
versioninfo:
	$(eval CURRENT_OS := $(shell uname -s))
	$(eval CURRENT_CPU_ARCH := $(shell uname -p))

	$(eval GIT_COMMIT_HASH := $(shell git rev-parse HEAD))
	$(eval GIT_TAG := $(shell git describe --tags --dirty))

	$(eval GO_LDFLAGS := -X ${MODULE_NAME}/internal/buildinfo.Version=${GIT_TAG} -X ${MODULE_NAME}/internal/buildinfo.GitCommitHash=${GIT_COMMIT_HASH})

# Display info about the current platform and configurations etc.
.PHONY: info
info: versioninfo
	@echo "Information"
	@echo "OS: ${CURRENT_OS}"
	@echo "CPU architecture: ${CURRENT_CPU_ARCH}"
	@echo "GIT_COMMIT_HASH: ${GIT_COMMIT_HASH}"
	@echo "GIT_TAG: ${GIT_TAG}"
	@echo "GO_LDFLAGS: ${GO_LDFLAGS}"

# Check that the source code is formatted correctly according to the gofmt standards
.PHONY: check-format
check-format:
	@echo "Checking formatting"
	@test -z $(shell gofmt -e -l ./ | tee /dev/stderr) || (echo "Please fix formatting first with gofmt" && exit 1)

# Check for other possible issues in the code
.PHONY: check-lint
check-lint:
	@echo "Linting code"
	go vet ./...
#NOTE: staticcheck is run as a github action and not as part of this Makefile

.PHONY: check
check: check-format check-lint
