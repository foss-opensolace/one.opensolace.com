GO := go
GORUN := $(GO) run

.PHONY: run

include .env
export

run: ./cmd/api/main.go
	$(GORUN) "$<"

psql:
	docker exec -it "$(POSTGRES_HOST)" psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DATABASE)" || \
	echo "Failed to connect!"