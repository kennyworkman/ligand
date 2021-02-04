NAME := ligand
RELEASE_DIR := release
ENVIRONMENT := dev

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LDFLAGS := -ldflags "-X github.com/latchai/ligand/go/pkg/global.Version=$(VERSION) -X github.com/latchai/ligand/go/pkg/global.Environment=$(ENVIRONMENT) -w"

BINARY = $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/$(NAME)
SHARED_BINARY = $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/ligand-daemon
MAIN := cmd/ligand/main.go
SHARED_MAIN := cmd/ligand-daemon/main.go
INSTALL_PATH := /usr/local/bin/$(NAME)

.PHONY: clean
clean:
	rm -rf $(RELEASE_DIR)

.PHONY: build
build: clean
	@mkdir -p $(RELEASE_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) $(MAIN)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(SHARED_BINARY) $(SHARED_MAIN)


# install without sudo if the install path exists and is writeable,
# or if it doesn't exist and its directory is writeable
# TODO: figure out if safe to always just `sudo` and simplify logic
.PHONY: install
install: build
	if [[ (-f "$(INSTALL_PATH)" && -w "$(INSTALL_PATH)") || (! -f "$(INSTALL_PATH)" && -w $$(dirname "$(INSTALL_PATH)")) ]]; then \
	    cp $(BINARY) $(INSTALL_PATH); \
	else \
	    sudo cp $(BINARY) $(INSTALL_PATH); \
	fi

# In addition to fixing imports, goimports also formats your code in the same
# style as gofmt so it can be used as a replacement for your editor's
# gofmt-on-save hook.
# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
#
.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports --local latchai -w -d .
