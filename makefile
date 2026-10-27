APP := writeit
BUILD_DIR := bin

.PHONY: build run test fmt vet clean install

build:
	go build -o $(BUILD_DIR)/$(APP) ./cmd

run:
	go run ./cmd

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)

install:
	go install ./cmd
