# Recursive wildcard.
# Usage: $(call rwildcard, dir_to_search, pattern_to_search_for)
# Example: $(call rwildcard, ., *.go)
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

# Go executable to use, i.e. `make GO=/usr/bin/go1.18`
# Defaults to first found in PATH
GO?=go

# TODO: download from latest release
bin/bindl:
	${GO} build -o bin/bindl -trimpath ./cmd/bindl

.PHONY: bin/bindl-dev
bin/bindl-dev: bin/goreleaser
	bin/goreleaser build \
		--output bin/bindl \
		--single-target \
		--snapshot \
		--rm-dist

include Makefile.*

.PHONY: archy
archy: bin/archy
	bin/archy -s -m

.PHONY: license
license: bin/addlicense
	bin/addlicense \
		-c "Bindl Authors" \
		-l apache \
		$(call rwildcard, ., *.go)

.PHONY: lint
lint: bin/golangci-lint
	bin/golangci-lint run

.PHONY: lint/fix
lint/fix: bin/golangci-lint
	bin/golangci-lint run --fix
