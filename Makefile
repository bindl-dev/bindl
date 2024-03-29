# Recursive wildcard.
# Usage: $(call rwildcard, dir_to_search, pattern_to_search_for)
# Example: $(call rwildcard, ., *.go)
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

# Go executable to use, i.e. `make GO=/usr/bin/go1.18`
# Defaults to first found in PATH
GO?=go


#########
# BUILD #
#########

bin/bindl:
	OUTDIR=bin ./bootstrap.sh

.PHONY: bin/bindl-min
bin/bindl-min:
	${GO} build -o bin/bindl ./cmd/bindl

.PHONY: bin/bindl-dev
bin/bindl-dev: bin/goreleaser
	bin/goreleaser build \
		--output bin/bindl \
		--single-target \
		--snapshot \
		--rm-dist

include Makefile.*

.PHONY: program/bootstrap/cosign-lock.yaml
program/bootstrap/cosign-lock.yaml: bin/bindl
	bin/bindl sync \
		--config program/bootstrap/cosign.yaml \
		--lock program/bootstrap/cosign-lock.yaml


###########
# RELEASE #
###########

.PHONY: release
release: bin/goreleaser bin/cosign bin/syft
	PATH=${PWD}/bin:${PATH} bin/goreleaser release --rm-dist


#########
# TESTS #
#########

program/testdata/myprogram.tar.gz:
	@./program/testdata/generate.sh

.PHONY: testdata
testdata: program/testdata/myprogram.tar.gz

.PHONY: test/unit
test/unit: testdata
	${GO} test -race -short -v ./...

.PHONY: test/integration
test/integration:
	${GO} test -race -run ".*[Ii]ntegration.*" -v ./...

# Manually build bindl and then download cosign because Makefile
# would not understand the dependency without bin/bindl existing.
.PHONY: test/functional
test/functional:
	PATH=${PWD}/bin:${PATH} ${MAKE} bin/bindl
	PATH=${PWD}/bin:${PATH} ${MAKE} bin/cosign
	PATH=${PWD}/bin:${PATH} ${GO} test -race -run ".*[Ff]unctional.*" -v ./...

.PHONY: test/all
test/all:
	${GO} test -race -v ./...


###########
# LINTERS #
###########

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

.PHONY: lint/gh-actions
lint/gh-actions: bin/golangci-lint
	bin/golangci-lint run --out-format github-actions 
