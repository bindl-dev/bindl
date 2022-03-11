include Makefile.*

.PHONY: bin/bindl-dev
bin/bindl-dev:
	go build -o bin/bindl ./cmd/bindl
