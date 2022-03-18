include Makefile.*

# Go executable to use, i.e. `make GO=/usr/bin/go1.18`
# Defaults to first found in PATH
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

.PHONY: license
license: bin/addlicense
	bin/addlicense \
		-c "Bindl Authors" \
		-l apache \
		**/*.go
