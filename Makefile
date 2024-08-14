# Go parameters
GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover
GOCOVREPORT=$(GOCOVER) -func=coverage.out

# Default target: Run tests
.PHONY: test
test:
	$(GOTEST) ./... -v

# Run tests with coverage
.PHONY: coverage
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCOVREPORT)

# Clean up coverage file
.PHONY: clean
clean:
	rm -f coverage.out

# Generate a test coverage report in HTML format
.PHONY: cover-html
cover-html:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCOVER) -html=coverage.out -o coverage.html

# Install dependencies
.PHONY: deps
deps:
	$(GO) mod tidy
	$(GO) mod download

