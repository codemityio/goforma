dep:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest graphviz --input-path=testdata/dep-pkg-recursive.dot --output-format=svg
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest graphviz --input-path=testdata/dep-pkg-recursive-exclude-standard.dot --output-format=svg

uml: uml-app uml-app-without-config uml-edge uml-recursive uml-recursive-uml-pkg-skip-const uml-recursive-uml-pkg-skip-func uml-recursive-uml-pkg-skip-notexported uml-recursive-uml-pkg-skip-primitive uml-recursive-uml-pkg-skip-var

uml-app:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-app.puml --output-format=svg --recursive

uml-app-without-config:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-app-without-config.puml --output-format=svg --recursive

uml-edge:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-edge.puml --output-format=svg --recursive

uml-recursive:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive.puml --output-format=svg --recursive

uml-recursive-uml-pkg-skip-const:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive-skip-const.puml --output-format=svg --recursive

uml-recursive-uml-pkg-skip-doc:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive-skip-doc.puml --output-format=svg --recursive

uml-recursive-uml-pkg-skip-func:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive-skip-func.puml --output-format=svg --recursive

uml-recursive-uml-pkg-skip-notexported:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive-skip-notexported.puml --output-format=svg --recursive

uml-recursive-uml-pkg-skip-primitive:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive-skip-primitive.puml --output-format=svg --recursive

uml-recursive-uml-pkg-skip-var:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/uml-pkg-recursive-skip-var.puml --output-format=svg --recursive
