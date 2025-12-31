build:
	@echo "Building..."
	@go build -o bin/one

run: build
	@echo "Running..."
	./bin/one
