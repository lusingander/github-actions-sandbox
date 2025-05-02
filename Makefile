.PHONY: build-mkdocs
build-mkdocs:
	mkdir -p docs/images
	cp -r assets/* docs/images/
	mkdocs build --strict

.PHONY: serve-mkdocs
serve-mkdocs:
	mkdir -p docs/images
	cp -r assets/* docs/images/
	mkdocs serve
