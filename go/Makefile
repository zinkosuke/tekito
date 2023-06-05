DIR := ./...
CI_TARGET := ci-lint ci-vet ci-errcheck ci-sqlc
FMT_TARGET := fmt-std fmt-tidy fmt-import

.PHONY: clean ci ci-* fmt fmt-* test test-* sqlc
clean: ${CI_TARGET}

ci: ${CI_TARGET}
ci-lint:
	staticcheck ${DIR}
ci-vet:
	go vet ${DIR}
ci-errcheck:
	# errcheck ${DIR}
ci-sqlc:
	sqlc diff > /dev/null 2>&1

fmt: ${FMT_TARGET}
fmt-std:
	go fmt ${DIR}
fmt-tidy:
	go mod tidy
fmt-import:
	find **/*.go -print | xargs goimports -w

test:
	go test -v ${DIR}

sqlc:
	sqlc generate
