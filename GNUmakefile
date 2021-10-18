PKG_NAME=barracudawaf

TEST?="./$(PKG_NAME)"
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
COVER_TEST?=$$(go list ./... |grep -v 'vendor')
WEBSITE_REPO=github.com/hashicorp/terraform-website

TAG := $(shell git describe --abbrev=0 --tags | cut -c2-)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
PLUGIN_NAME := terraform-provider-$(PKG_NAME)
TF_PLUGIN_PATH := $(HOME)/.terraform.d/plugins/registry.terraform.io/hashicorp/$(PKG_NAME)/$(TAG)/$(GOOS)_$(GOARCH)

default: build

build: fmtcheck
	go get -d
	go build -o $(PLUGIN_NAME)

install: build
	install -d $(TF_PLUGIN_PATH) && mv $(PLUGIN_NAME) $(TF_PLUGIN_PATH)

plugin: 
	install -d $(TF_PLUGIN_PATH) && mv $(PLUGIN_NAME) $(TF_PLUGIN_PATH)

test: fmtcheck
	go test $(TEST) -v || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

testrace: fmtcheck
	TF_ACC= go test -race $(TEST) $(TESTARGS)

compile: fmtcheck
	@sh -c "'$(CURDIR)/scripts/compile.sh'"

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(COVER_TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile: fmtcheck
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./barracudawaf"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build install plugin test testacc testrace cover vet fmt fmtcheck errcheck test-compile
