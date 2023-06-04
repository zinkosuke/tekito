TARGET := ./...
cmd_poetry := poetry run

.PHONY: clean ci-* test
clean: ci-fmt ci-vet ci-goimports ci-errcheck ci-lint
ci-fmt:
	go fmt ${TARGET}
ci-vet:
	go vet ${TARGET}
ci-lint:
	staticcheck ${TARGET}
ci-goimports:
	find **/*.go -print | xargs goimports -w
	# find **/*.go -print | xargs goimports -w -local ""
ci-errcheck:
	# errcheck ${TARGET}
	echo
test:
	go test -v ${TARGET}
