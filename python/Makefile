
#
# Testing.
#

.PHONY: test
test: lint
	pytest --import-mode=importlib -m "not external"

.PHONY: test-external
test-external: lint
	pytest --import-mode=importlib


#
# Code Hygiene.
#

.PHONY: fmt
fmt:
	black .

.PHONY: lint
lint:
	mypy .
	black --check .
