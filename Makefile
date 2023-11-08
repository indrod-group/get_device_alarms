.PHONY: build run clean

BINARY_NAME=alarms_notification

build:
	@echo "Building..."
	go build -o ./bin/$(BINARY_NAME) -ldflags "-s -w"
	@echo "Build complete"

run: build
	@echo "Running..."
	./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	go clean
	rm ./bin/$(BINARY_NAME)
	@echo "Clean complete"
