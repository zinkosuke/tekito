DIR := ./...

.PHONY: ci ci-*
ci: ci-lint ci-vet ci-errcheck
ci-lint:
	staticcheck ${DIR}
ci-vet:
	go vet ${DIR}
ci-errcheck:
	# errcheck ${DIR}

.PHONY: fmt fmt-*
fmt: fmt-std fmt-tidy fmt-import
fmt-std:
	go fmt ${DIR}
fmt-tidy:
	go mod tidy
fmt-import:
	find -type f -name '*.go' -print | xargs goimports -w

.PHONY: test
test:
	go test -v ${DIR}
