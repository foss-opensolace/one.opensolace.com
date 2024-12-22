GO := go
GORUN := $(GO) run

.PHONY: run

run: ./cmd/api/main.go
	$(GORUN) "$<"