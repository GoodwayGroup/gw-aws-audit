# Build Variables
NAME     = gw-aws-audit
VERSION ?= $(shell git describe --tags --always)

# Go variables
GO      ?= go
GOOS    ?= $(shell $(GO) env GOOS)
GOARCH  ?= $(shell $(GO) env GOARCH)
GOHOST  ?= GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO)

LDFLAGS ?= "-X main.version=$(VERSION)"

.PHONY: all
all: help

###############
##@ Development

.PHONY: clean
clean: ## Clean workspace
	@ $(MAKE) --no-print-directory log-$@
	@rm -rf bin/ && rm -rf build/ && rm -rf dist/ && rm -rf cover.out

.PHONY: test
test: ## Run tests
	@ $(MAKE) --no-print-directory log-$@
	$(GOHOST) test -covermode atomic -coverprofile cover.out -v ./...

.PHONY: lint
lint:   ## Run linters
	@ $(MAKE) --no-print-directory log-$@
	golangci-lint run

#########
##@ Build

.PHONY: build
build: clean ## Build gw-aws-audit
	@ $(MAKE) --no-print-directory log-$@
	@mkdir -p bin/
	CGO_ENABLED=0 $(GOHOST) build -ldflags=$(LDFLAGS) -o bin/$(NAME) ./main.go

###########
##@ Release

.PHONY: changelog
changelog: ## Generate changelog
	@ $(MAKE) --no-print-directory log-$@
	git-chglog --next-tag $(VERSION) -o CHANGELOG.md

.PHONY: release
release: ## Release a new tag
	@ $(MAKE) --no-print-directory log-$@
	./release.sh $(VERSION)

.PHONY: docs
docs: ## Generate new docs
	@ $(MAKE) --no-print-directory log-$@
	DOCS_MD=1 go run ./main.go > docs/$(NAME).md
	DOCS_MAN=1 go run ./main.go > docs/$(NAME).8

########
##@ Help

.PHONY: help
help:   ## Display this help
	@awk \
		-v "col=\033[36m" -v "nocol=\033[0m" \
		' \
			BEGIN { \
				FS = ":.*##" ; \
				printf "Usage:\n  make %s<target>%s\n", col, nocol \
			} \
			/^[a-zA-Z_-]+:.*?##/ { \
				printf "  %s%-12s%s %s\n", col, $$1, nocol, $$2 \
			} \
			/^##@/ { \
				printf "\n%s%s%s\n", nocol, substr($$0, 5), nocol \
			} \
		' $(MAKEFILE_LIST)

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk \
			'BEGIN { \
				FS = ":.*?## " \
			}; \
			{ \
				printf "\033[36m==> %s\033[0m\n", $$2 \
			}'
