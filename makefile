APP := writeit
BUILD_DIR := bin

.PHONY: build run test fmt vet clean install

build:
	reset
	go build -o $(BUILD_DIR)/$(APP) ./cmd

run:
	reset
	go run ./cmd

test:
	reset
	go test ./...

fmt:
	reset
	go fmt ./...

vet:
	reset
	go vet ./...

clean:
	reset
	rm -rf $(BUILD_DIR)

install:
	reset
	go install ./cmd
