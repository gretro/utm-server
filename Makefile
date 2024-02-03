OUT_DIR="bin"

build: bin/utm_server_amd64 bin/utm_server_arm64
	@echo "Build complete"
.PHONY: build

clean:
	@echo "Cleaning up"
	@rm -rf $(OUT_DIR)
.PHONY: clean

bin/utm_server_amd64:
	@echo "Building for Darwin AMD64"
	@GOARCH=amd64 GOOS=darwin go build -o $(OUT_DIR)/utm_server_amd64 ./src/cmd/api.go

bin/utm_server_arm64:
	@echo "Building for Darwin ARM64"
	@GOARCH=arm64 GOOS=darwin go build -o $(OUT_DIR)/utm_server_arm64 ./src/cmd/api.go

dev:
	@echo "Starting development server"
	@go run ./src/cmd/api.go