# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

.PHONY: clean

clean:
	$(GOCLEAN)

test:
	$(GOTEST) ./...
