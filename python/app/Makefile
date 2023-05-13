TARGET := .
cmd_poetry := poetry run

.PHONY: clean fmt-* lint-* test
clean: fmt-full lint-full
fmt-full: fmt-isort fmt-black
fmt-isort:
	$(cmd_poetry) isort ${TARGET}
fmt-black:
	$(cmd_poetry) black ${TARGET}
lint-full: lint-isort lint-flake8 lint-mypy
lint-isort:
	$(cmd_poetry) isort --check-only ${TARGET}
lint-flake8:
	$(cmd_poetry) pflake8 ${TARGET}
lint-mypy:
	$(cmd_poetry) mypy --install-types --non-interactive
	$(cmd_poetry) mypy ${TARGET}
test:
	$(cmd_poetry) pytest ${TARGET}
