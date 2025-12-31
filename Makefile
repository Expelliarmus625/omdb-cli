build:
	@echo "Building..."
	@go build -o bin/omdb-cli

run: build
	@echo "Running..."
	./bin/omdb-cli
