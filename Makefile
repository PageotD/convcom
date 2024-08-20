# Variables
GO_CMD = go
SRC_DIR = cmd/convcom
BIN_DIR = bin
INT_DIR = internal
INT_PKG = $(wildcard $(INT_DIR)/*)
BIN_NAME = convcom
PRF_DIR = profiling
BUILD_FLAGS = -ldflags="-s -w" -a -tags netgo -installsuffix netgo
VERSION = $(shell cat version.txt)

# Targets
all: build

# Compile for specified OS and ARCH
build-for-os-arch:
	mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_CMD) build -C $(SRC_DIR) -o ../../$(BIN_DIR)/$(BIN_NAME) $(BUILD_FLAGS)

# Compile for all platforms
build:
	$(MAKE) build-for-os-arch GOOS=linux GOARCH=amd64

# Test the cmd/taskido package
test-cmd:
	$(GO_CMD) -C $(SRC_DIR) test -v -coverprofile=../../$(PRF_DIR)/coverage-cmd.out -cpuprofile=../../$(PRF_DIR)/cpu-cmd.out -memprofile=../../$(PRF_DIR)/mem-cmd.out 

# Test the internal package
test-internal:
	for pkg in $(INT_PKG); do\
		$(GO_CMD) -C $${pkg} test -v -coverprofile=../../$(PRF_DIR)/coverage-internal.out -cpuprofile=../../$(PRF_DIR)/cpu-internal.out -memprofile=../../$(PRF_DIR)/mem-internal.out;\
	done

# Run all tests
test: test-cmd test-internal

# Clean build artifacts
clean:
	rm -f $(BIN_DIR)/*

.PHONY: all build clean test test-cmd test-internal
