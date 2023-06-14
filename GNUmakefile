NAME=kamatera
BINARY=packer-plugin-${NAME}

COUNT?=1
TEST?=$(shell go list ./...)
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

.PHONY: dev

build:
	@go build -o ${BINARY}

generate:
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@latest
	@go generate -v ./...

ci-release-docs:
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@latest
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst docs/
	@/bin/sh -c "[ -d docs ] && zip -r docs.zip docs/"

dev: build
	@mkdir -p ~/.packer.d/plugins/
	@mv ${BINARY} ~/.packer.d/plugins/${BINARY}

run-example: dev
	@packer build ./example

test:
	@go test -count $(COUNT) $(TEST) -timeout=3m

install-packer-sdc: ## Install packer sofware development command
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

build-docs: install-packer-sdc
	@if [ -d ".docs" ]; then rm -r ".docs"; fi
	@packer-sdc renderdocs -src "docs" -partials docs-partials/ -dst ".docs/"
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "BrandonRomano"
	@rm -r ".docs"
