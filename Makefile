BINARY_NAME=netforge
BUILD_DIR=bin
VERSION=1.0.1

.PHONY: all build windows linux darwin clean version

all: build

version:
	@echo "NetForge version $(VERSION)"

build: windows windows386 windowsarm64 linux darwin

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

windows386:
	GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe main.go

windowsarm64:
	GOOS=windows GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

clean:
	rm -rf $(BUILD_DIR)
