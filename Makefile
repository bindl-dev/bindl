include Makefile.*

.PHONY: bin/bindl-dev
bin/bindl-dev:
	go build -o bin/bindl -trimpath ./cmd/bindl

# TODO: download from latest release
bin/bindl:
	go build -o bin/bindl -trimpath ./cmd/bindl

bin/archy: bin/bindl
	bin/bindl get archy

bin/golangci-lint: bin/bindl
	bin/bindl get golangci-lint

# doesn't actually work for 1.18 yet: https://github.com/go-critic/go-critic/issues/1157
.PHONY: lint
lint: bin/golangci-lint
	bin/golangci-lint run
