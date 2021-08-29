.DEFAULT_GOAL := help
.PHONY: deps install-deps-no-envs install-docker-deps install-deps install-linters
.PHONY: install-deps-Linux install-deps-Darwin install-deps-Windows
.PHONY: build build-docker build-icon
.PHONY: prepare-release gen-mocks
.PHONY: run help
.PHONY: test test-core test-sky test-sky-launch-html-cover test-cover lint
.PHONY: clean-test clean-build clean clean-Windows

# Application info (for dumping)
ORG_DOMAIN		:= simelo.tech.org
ORG_NAME		:= Simelo.Tech
APP_NAME		:= FiberCryptoWallet
APP_DESCRIPTION	:= Multi-coin cryptocurrency wallet
APP_VERSION		:= 0.27.0
LICENSE			:= GPLv3
COPYRIGHT		:= Copyright © 2019 $(ORG_NAME)

UNAME_S = $(shell uname -s)
OSNAME = $(shell echo $(UNAME_S) | tr A-Z a-z)
DEFAULT_TARGET ?= desktop
DEFAULT_ARCH ?= linux
## In future use as a parameter tu make command.
COIN ?= skycoin
COVERAGEPATH = src/coin/$(COIN)
COVERAGEPREFIX = $(COVERAGEPATH)/coverage
COVERAGEFILE = $(COVERAGEPREFIX).out
COVERAGETEMP = $(COVERAGEPREFIX).tmp.out
COVERAGEHTML = $(COVERAGEPREFIX).html

# Icons
APP_ICON_PATH	:= resources/images/icons/appIcon
ICONS_BUILDPATH	:= resources/images/icons/appIcon/build
ICONSET			:= resources/images/icons/appIcon/appIcon.iconset
CONVERT			:= convert
SIPS			:= sips
ICONUTIL		:= iconutil
UNAME_S         = $(shell uname -s)
DEFAULT_TARGET  ?= desktop
DEFAULT_ARCH    ?= linux

# Platform-specific switches
ifeq ($(OS),Windows_NT)
	CONVERT		 = ./convert.exe
	WINDRES		:= windres
	RC_FILE		:= resources/platform/windows/winResources.rc
	RC_OBJ		:= winResources.syso
else
#	Get the UNIX operating system
	OS			:= $(shell uname -s)
#	If Linux...
	ifeq ($(OS),Linux)
	endif
#	If Darwin
	ifeq ($(OS),Darwin)
		DARWIN_RES	:= darwin
		PLIST		:= resources/platform/darwin/info.plist
	endif
endif

# Files
GIT := $(shell which git)
ifeq ($(GIT),)
  ALLFILES    := $(shell find . type f | grep -v .git | grep -v vendor)
else
	GIT_BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)
	ALLFILES    := $(shell git ls-tree -r $(GIT_BRANCH) --name-only | grep -v vendor)
endif

GOFILES    := $(shell echo "$(ALLFILES)" | grep  '.go$$')
QRCFILES   := $(shell echo "$(ALLFILES)" | grep  '.qrc$$')
QMLFILES   := $(shell echo "$(ALLFILES)" | grep  '.qml$$')
TSFILES    := $(shell echo "$(ALLFILES)" | grep  '.ts$$')
SVGFILES   := $(shell echo "$(ALLFILES)" | grep  '.svg$$')
JSFILES    := $(shell echo "$(ALLFILES)" | grep  '.js$$')
PNGFILES   := $(shell echo "$(ALLFILES)" | grep  '.png$$')
OTFFILES   := $(shell echo "$(ALLFILES)" | grep  '.otf$$')
ICNSFILES  := $(shell echo "$(ALLFILES)" | grep  '.icns$$')
ICOFILES   := $(shell echo "$(ALLFILES)" | grep  '.ico$$')
RCFILES    := $(shell echo "$(ALLFILES)" | grep  '.rc$$')
PLISTFILES := $(shell echo "$(ALLFILES)" | grep  'Info.plist$$')
QTCONFFILES := $(shell echo "$(ALLFILES)" | grep  'qtquickcontrols2.conf$$')

QMLFILES      := $(shell echo "$(QMLFILES) $(JSFILES)")
QTFILES       := $(shell echo "$(QRCFILES) $(TSFILES) $(PLISTFILES) $(QTCONFFILES)")
RESOURCEFILES := $(shell echo "$(SVGFILES) $(PNGFILES) $(OTFFILES) $(ICNSFILES) $(ICOFILES) $(RCFILES)")
SRCFILES      := $(shell echo "$(QTFILES) $(RESOURCEFILES) $(GOFILES)")

BINPATH_Linux      := deploy/linux/fibercryptowallet
BINPATH_Windows_NT := deploy/windows/fibercryptowallet.exe
BINPATH_Darwin     := deploy/darwin/fibercryptowallet.app/Contents/MacOS/fibercryptowallet
BINPATH            := $(BINPATH_$(UNAME_S))

PWD := $(shell pwd)

GOPATH ?= $(shell echo "$${GOPATH}")
GOPATH_SRC := src/github.com/fibercrypto/fibercryptowallet

DOCKER_QT       ?= therecipe/qt
DOCKER_QT_TEST  ?= simelotech/qt-test

deps: ## Add dependencies
	dep ensure

# Targets
install-deps-no-envs: ## Install therecipe/qt with -tags=no_env set
	go get -v -tags=no_env github.com/therecipe/qt/cmd/...
	go get -t -d -v ./...
	@echo "Dependencies installed"

install-docker-deps: ## Install docker images for project compilation using docker
	@echo "Downloading images..."
	docker pull $(DOCKER_QT):$(DEFAULT_ARCH)
	docker pull $(DOCKER_QT_TEST):$(DEFAULT_ARCH)
	@echo "Download finished."

install-deps-Linux: ## Install Linux dependencies
	sudo apt-get update
	go get -t -d -v ./...

install-deps-Darwin: ## Install osx dependencies
	xcode-select --install || true
	go get -t -d -v ./...

install-deps-Windows: ## Install Windowns dependencies
	set GO111MODULE=off
	go get -t -d -v ./...

install-deps: install-deps-$(UNAME_S) install-linters ## Install dependencies
	@echo "Dependencies installed"

build-docker: install-docker-deps ## Build project using docker
	@echo "Building $(APP_NAME)..."
	@echo "No build rules defined"
	@echo "Done."

build: $(LIBPATH)  ## Build FiberCrypto Wallet
	@echo "Output => $(BINPATH)"

prepare-release: ## Change the resources in the app and prepare to release the app
	./.travis/setup_release.sh

clean-test: ## Remove temporary test files
	rm -f $(COVERAGEFILE)
	rm -f $(COVERAGETEMP)
	rm -f $(COVERAGEHTML)

clean-build: ## Remove temporary files
	@echo "Cleaning project $(APP_NAME)..."
	rm -rf deploy/
	rm -rf linux/
	rm -rf windows/
	rm -rf rcc.cpp
	rm -rf rcc.qrc
	rm -rf rcc_cgo_*.go
	rm -rf rcc_*.cpp
	rm -rf rcc__*
	find . -path "*moc.*" -delete
	find . -path "*moc_*" -delete
	rm -rf "$(ICONS_BUILDPATH)"
	rm -rf "$(RC_OBJ)"
	rm -rf "$(ICONSET)"

	@echo "Done."

clean: clean-test clean-build ## Remove temporary files

gen-mocks-core: ## Generate mocks for core interface types
	mockery -all -output src/coin/mocks -outpkg mocks -dir src/core

gen-mocks-sky: ## Generate mocks for sky-wallet interface types
	mockery -name Devicer -dir ./vendor/github.com/fibercrypto/skywallet-go/src/skywallet -output ./src/contrib/skywallet/mocks -case underscore
	mockery -name DeviceDriver -dir ./vendor/github.com/fibercrypto/skywallet-go/src/skywallet -output ./src/contrib/skywallet/mocks -case underscore

gen-mocks: gen-mocks-core gen-mocks-sky ## Generate mocks for interface types

$(COVERAGEFILE):
	echo 'mode: set' > $(COVERAGEFILE)

test-skyhw: ## Run Hardware wallet tests
	go test -coverprofile=$(COVERAGETEMP) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/contrib/skywallet
	cat $(COVERAGETEMP) | grep -v '^mode: set$$' >> $(COVERAGEFILE)

test-sky: ## Run Skycoin plugin test suite
	go test -coverprofile=$(COVERAGETEMP) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/coin/skycoin
	cat $(COVERAGETEMP) | grep -v '^mode: set$$' >> $(COVERAGEFILE)
	go test -coverprofile=$(COVERAGETEMP) -timeout 60s github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/models
	cat $(COVERAGETEMP) | grep -v '^mode: set$$' >> $(COVERAGEFILE)

test-core: ## Run tests for API core and helpers
	go test -coverprofile=$(COVERAGETEMP) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/util
	cat $(COVERAGETEMP) | grep -v '^mode: set$$' >> $(COVERAGEFILE)

test-data: ## Run tests for data package
	go test -coverprofile=$(COVERAGETEMP) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/data
	cat $(COVERAGETEMP) | grep -v '^mode: set$$' >> $(COVERAGEFILE)

test-html-cover:
	go tool cover -html=$(COVERAGEFILE) -o $(COVERAGEPREFIX).html

test-cover-travis: clean-test
	go test -covermode=count -coverprofile=$(COVERAGEFILE) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/util
	$(GOPATH)/bin/goveralls -coverprofile=$(COVERAGEFILE) -service=travis-ci -repotoken 1zkcSxi8TkcxpL2zTQOK9G5FFoVgWjceP
	go test -coverprofile=$(COVERAGEFILE) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/models
	$(GOPATH)/bin/goveralls -coverprofile=$(COVERAGEFILE) -service=travis-ci -repotoken 1zkcSxi8TkcxpL2zTQOK9G5FFoVgWjceP
	go test -cover -covermode=count -coverprofile=$(COVERAGEFILE) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/coin/skycoin
	$(GOPATH)/bin/goveralls -coverprofile=$(COVERAGEFILE) -service=travis-ci -repotoken 1zkcSxi8TkcxpL2zTQOK9G5FFoVgWjceP
	go test -cover -covermode=count -coverprofile=$(COVERAGEFILE) -timeout 30s github.com/fibercrypto/fibercryptowallet/src/contrib/skywallet
	$(GOPATH)/bin/goveralls -coverprofile=$(COVERAGEFILE) -service=travis-ci -repotoken 1zkcSxi8TkcxpL2zTQOK9G5FFoVgWjceP

test-cover: test test-html-cover ## Show more details of test coverage

test: clean-test $(COVERAGEFILE) test-core test-sky test-data ## Run project test suite

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	cat ./.travis/install-golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.21.0

install-coveralls: ## Install coveralls
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

lint: ## Run linters. Use make install-linters first.
	# src needs separate linting rules
	golangci-lint run -c .golangci.yml ./src/coin/...
	golangci-lint run -c .golangci.yml ./src/core/...
	golangci-lint run -c .golangci.yml ./src/main/...
	golangci-lint run -c .golangci.yml ./src/util/...

help:
	@echo "$(APP_NAME) v$(APP_VERSION)"
	@echo "$(APP_DESCRIPTION)"
	@echo "$(COPYRIGHT)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
