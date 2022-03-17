include Makefile.*

GO?=go

.PHONY: bin/bindl-dev
bin/bindl-dev:
	${GO} build -o bin/bindl -trimpath ./cmd/bindl

# TODO: download from latest release
bin/bindl:
	${GO} build -o bin/bindl -trimpath ./cmd/bindl

.PHONY: archy
archy: bin/archy
	bin/archy -s -m
