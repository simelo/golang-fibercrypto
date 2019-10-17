.DEFAULT_GOAL := help
.PHONY: run build clean help

run:  ## Run FiberCrypto Wallet.
	@echo "Running FiberCrypto Wallet..."
	@./deploy/linux/FiberCryptoWallet

build:  ## Build FiberCrypto Wallet.
	@echo "Building FiberCrypto Wallet..."
	# Add the flag `-quickcompiler` when making a release
	@qtdeploy build desktop
	@echo "Done."

mocks: ## Create all mock files for unit tests
	mockery -name Devicer -dir ./vendor/github.com/skycoin/hardware-wallet-go/src/skywallet -output ./src/hardware/mocks -case underscore
	mockery -name DeviceDriver -dir ./vendor/github.com/skycoin/hardware-wallet-go/src/skywallet -output ./src/hardware/mocks -case underscore


clean: ## Clean project FiberCrypto Wallet.
	@echo "Cleaning project FiberCrypto Wallet..."
	rm -rf deploy/
	rm -rf linux/
	rm -rf rcc.cpp
	rm -rf rcc.qrc
	rm -rf rcc_cgo_linux_linux_amd64.go
	rm -rf rcc_*.cpp
	rm -rf rcc__*
	find . -path "*moc.*" -delete
	find . -path "*moc_*" -delete
	@echo "Done."

test-hw: mocks ## Run Hardware wallet tests
	go test github.com/fibercrypto/FiberCryptoWallet/src/hardware

test-sky: ## Run Skycoin plugin test suite
	go test github.com/fibercrypto/FiberCryptoWallet/src/coin/skycoin
	go test github.com/fibercrypto/FiberCryptoWallet/src/coin/skycoin/models

test: test-sky test-hw ## Run project test suite

test: test-sky ## Run project test suite

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
