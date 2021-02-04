.DEFAULT_GOAL := help
CLI?=latch
DAEMON?=latch-daemon
COMMIT_SHA?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

clean:
		@echo "$(OK_COLOR)==> Cleaning working directory... $(NO_COLOR)"
		@go clean
		rm -f $(CLI)

build: clean
		@echo "$(OK_COLOR)==> Building binary... $(NO_COLOR)"
		go build -o $(CLI) cmd/$(CLI)/main.go
		go build -o $(DAEMON) cmd/$(DAEMON)/main.go
		mkdir -p ../python/latch/bin
		cp latch-daemon ../python/latch/bin

run:
		go run cmd/$(DAEMON)/main.go '/var/folders/ch/kcj80q7j7v1g99sb4ddslmp80000gn/T/latch-daemon-nihp5mef.sock'


# In addition to fixing imports, goimports also formats your code in the same
# style as gofmt so it can be used as a replacement for your editor's
# gofmt-on-save hook.
# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
#
.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports --local latchai -w -d .
