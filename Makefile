DIST := dist
IMPORT := code.gitea.io/tea
export GO111MODULE=on

GO ?= go
SED_INPLACE := sed -i
SHASUM ?= shasum -a 256

export PATH := $($(GO) env GOPATH)/bin:$(PATH)

ifeq ($(OS), Windows_NT)
	EXECUTABLE := tea.exe
else
	EXECUTABLE := tea
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		SED_INPLACE := sed -i ''
	endif
endif

GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
GOFMT ?= gofmt -s

GOFLAGS := -i -v
EXTRA_GOFLAGS ?=

MAKE_VERSION := $(shell make -v | head -n 1)

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
	TEA_VERSION ?= $(VERSION)
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
	TEA_VERSION ?= $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
endif

TEA_VERSION_TAG ?= $(shell sed 's/+/_/' <<< $(TEA_VERSION))

LDFLAGS := -X "main.Version=$(TEA_VERSION)" -X "main.Tags=$(TAGS)"

PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)

TAGS ?=

ifeq ($(OS), Windows_NT)
	EXECUTABLE := tea.exe
else
	EXECUTABLE := tea
endif

# $(call strip-suffix,filename)
strip-suffix = $(firstword $(subst ., ,$(1)))

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -mod=vendor -i ./...
	rm -rf $(EXECUTABLE) $(DIST)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: vet
vet:
	# Default vet
	$(GO) vet -mod=vendor $(PACKAGES)
	# Custom vet
	$(GO) build -mod=vendor code.gitea.io/gitea-vet
	$(GO) vet -vettool=gitea-vet $(PACKAGES)

.PHONY: lint
lint:
	@hash revive > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/mgechev/revive; \
	fi
	revive -config .revive.toml -exclude=./vendor/... ./... || exit 1

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error -i unknwon,destory $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w -i unknwon $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: test
test:
	$(GO) test -mod=vendor -tags='sqlite sqlite_unlock_notify' $(PACKAGES)

.PHONY: unit-test-coverage
unit-test-coverage:
	$(GO) test -mod=vendor -tags='sqlite sqlite_unlock_notify' -cover -coverprofile coverage.out $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1

.PHONY: vendor
vendor:
	$(GO) mod tidy && $(GO) mod vendor

.PHONY: test-vendor
test-vendor: vendor
	@diff=$$(git diff vendor/); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make vendor' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: check
check: test

.PHONY: install
install: $(wildcard *.go)
	$(GO) install -mod=vendor -v -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)'

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	$(GO) build -mod=vendor $(GOFLAGS) $(EXTRA_GOFLAGS) -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)' -o $@

.PHONY: build-image
build-image:
	docker build --build-arg VERSION=$(TEA_VERSION) -t gitea/tea:$(TEA_VERSION_TAG) .

.PHONY: release
release: release-dirs release-os release-compress release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-os
release-os:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/mitchellh/gox; \
	fi
	CGO_ENABLED=0 gox -verbose -cgo=false -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)' -osarch='!darwin/386 !darwin/arm64 !darwin/arm' -os="windows linux darwin" -arch="386 amd64 arm arm64" -output="$(DIST)/release/tea-$(VERSION)-{{.OS}}-{{.Arch}}"

.PHONY: release-compress
release-compress:
	@hash gxz > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		GO111MODULE=off $(GO) get -u github.com/ulikunitz/xz/cmd/gxz; \
	fi
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "compressing $${file}" && gxz -k -9 $${file}; done;

.PHONY: release-check
release-check:
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "checksumming $${file}" && $(SHASUM) `echo $${file} | sed 's/^..//'` > $${file}.sha256; done;
