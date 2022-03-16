include Makefile.*

.PHONY: bin/bindl-dev
bin/bindl-dev:
	go build -o bin/bindl -trimpath ./cmd/bindl

# TODO: download from latest release
bin/bindl:
	go build -o bin/bindl -trimpath ./cmd/bindl

.PHONY: archy
archy: bin/archy
	bin/archy -s -m
