# Temp/restic-browser/Makefile

# Variables
BINARY=restic-browser
SRC=main.go
PORT?=8081

.PHONY: all build clean run help

all: build

build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) $(SRC)
	@echo "Build complete!"

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY)

run: build
	@echo "Usage: make run REPO=/path/to/repo PASSWORD=yourpassword [PORT=8081]"
	@if [ -z "$(REPO)" ] || [ -z "$(PASSWORD)" ]; then \
		echo "Error: REPO and PASSWORD variables must be set. Example:"; \
		echo "  make run REPO=/path/to/repo PASSWORD=yourpassword [PORT=8081]"; \
		exit 1; \
	fi
	@echo "Starting $(BINARY) on port $(PORT) for repo $(REPO)..."
	./$(BINARY) -repo "$(REPO)" -password "$(PASSWORD)" -port "$(PORT)"

help:
	@echo "Restic Browser Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build the Go binary"
	@echo "  clean         Remove the built binary"
	@echo "  run           Build and run the service (requires REPO and PASSWORD variables)"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run REPO=/path/to/repo PASSWORD=yourpassword PORT=9000"
