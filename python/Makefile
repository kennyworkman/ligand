
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

.PHONY: clean
clean:
	rm -rf build dist replicate/bin

.PHONY: build
build: clean
	pip install wheel
	python setup.py bdist_wheel --plat-name manylinux1_x86_64
	python setup.py bdist_wheel --plat-name macosx_10_9_x86_64
