MAIN_FILE = ./cmd/server/main.go
BUILD_FILE = ./go.learning

build:
	go build -o $(BUILD_FILE) $(MAIN_FILE)

run:
	go run $(MAIN_FILE)

run-build:
	$(BUILD_FILE)

clean:
	rm $(BUILD_FILE)

.PHONY: run build run-build clean
