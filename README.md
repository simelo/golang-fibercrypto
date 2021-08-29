# FiberCrypto wallet

[![Build Status](https://travis-ci.org/fibercrypto/golang-fibercrypto.svg?branch=develop)](https://travis-ci.org/fibercrypto/golang-fibercrypto)
[![Contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE.GPLv3)
[![Coverage Status](https://coveralls.io/repos/github/fibercrypto/FiberCryptoWallet/badge.svg?branch=develop)](https://coveralls.io/github/fibercrypto/FiberCryptoWallet?branch=develop)

Welcome to the FiberCrypto API project repository. The source code tracked in this repository serves as a backend to the [FiberCrypto v2 wallet software](https://github.com/fibercrypto/fibercrypto). The goals we work towards are the following:

- Define and implement core system architecture
- Out-of-the-box support for every SkyFiber token in a single place
- Support [Bitcoin](http://bitcoin.org) and other altcoins
- Implement integrations for crypto assets exchange
- Integrations with trading tools
- Data sources for basic blockchain-specific tools

## Development

### Project folder structure

Project files are organized as follows:

- `main.go` : Library entry point(s)
- `CHANGELOG.md` : Project changelog
- `Makefile` : Project build rules
- `README.md` : This file.
- `./src` : Librarry source code.
- `./src/core` : Core go-lang interfaces.
- `./src/main` : Project specific source code.
- `./src/util` : Reusable code.
- `./src/util/logging` : Event logging infrastructure.
- `./src/coin` : Source code for altcoin integrations.
- `./src/coin/mocks` : Types implementing `core` interfaces for generic testing scenarios
- `./src/coin/skycoin` : Skycoin wallet integration
- `./src/coin/skycoin/models` : Skycoin implementation of golang core interfaces.
- `./src/coin/skycoin/blockchain` : Skycoin blockchain API.
- `./src/coin/skycoin/sign` : Skycoin sign API.
- `./src/contrib` : Extra extensions, plugins and tools
- `./src/contrib/skywallet/` : SkyWallet sign API implementation
- `./vendor` : Project dependencies managed by `dep`.

### Architecture

FiberCrypto core supports multiple altcoins. In order to cope with this complexity wallet GUI code, QT models, DApps, libraries and the whole FiberCrypto ecosystem rely on strict interfaces which shall be implemented to add support for a given crypto asset coin. Each such integration must have the following main components:

- `Models API`: Implements application core interfaces.
- `Sign API` : Implements altcoin transaction and message signing primitives required by application code.
- `Blockchain API` : Provides communication between application and altcoin service nodes to query for data via REST, JSON-RPC and other similar low-level client-server API.
- `Peer-exchange API` (optional): Implements peer-to-peer interactions with altcoin blockchain nodes.

### Build System

The build system is the standard [Go](https://golang.org/ "The Go Programming Language") build machhinery.

#### Requirements

No platform-soecific requirements identified so far

#### Make targets

Common actions are automated with the help of `make`. The following targets have been implemnented:

```
deps                           Add dependencies
install-deps-no-envs           Install therecipe/qt with -tags=no_env set
install-docker-deps            Install docker images for project compilation using docker
install-deps-Linux             Install Linux dependencies
install-deps-Darwin            Install osx dependencies
install-deps-Windows           Install Windowns dependencies
install-deps                   Install dependencies
build-docker                   Build project using docker
build                          Build FiberCrypto API
prepare-release                Change the resources in the app and prepare to release the app
clean-test                     Remove temporary test files
clean-build                    Remove temporary files
clean                          Remove temporary files
gen-mocks-core                 Generate mocks for core interface types
gen-mocks-sky                  Generate mocks for internal Skycoin types
gen-mocks                      Generate mocks for interface types
test-sky                       Run Skycoin plugin test suite
test-core                      Run tests for API core and helpers
test-data                      Run tests for data package
test-cover                     Show more details of test coverage
test                           Run project test suite
install-linters                Install linters
install-coveralls              Install coveralls
lint                           Run linters. Use make install-linters first.
```

Type `make help` in your console for details.

## Releases

### Update the version

0. If the `master` branch has commits that are not in `develop` (e.g. due to a hotfix applied to `master`), merge `master` into `develop`
0. Update `CHANGELOG.md`: move the "unreleased" changes to the version and add the date
0. Update the files in https://github.com/skycoin/repo-info by following the [metadata update procedure](https://github.com/skycoin/repo-info/#updating-skycoin-repository-metadate),
0. Merge these changes to `develop`
0. Follow the steps in [pre-release testing](#pre-release-testing)
0. Make a PR merging `develop` into `master`
0. Review the PR and merge it
0. Tag the `master` branch with the version number. Version tags start with `v`, e.g. `v0.1.0`.
    Sign the tag. If you have your GPG key in github, creating a release on the Github website will automatically tag the release.
    It can be tagged from the command line with `git tag -as v0.20.0 $COMMIT_ID`, but Github will not recognize it as a "release".
0. Make sure that the app runs properly from the `master` branch
0. Release builds are created and uploaded by travis. To do it manually, checkout the `master` branch and follow the [create release builds](#creating-release-builds) instructions.

If there are problems discovered after merging to `master`, start over, and increment the 3rd version number.
For example, `v0.1.0` becomes `v0.1.1`, for minor fixes.

### Pre-release testing

Performs these actions before releasing:

* `make test-sky` Run Skycoin plugin test suite
* `make test-core` Run tests for API core and helpers
* `make test-data` Run tests for data package
* `make test-cover` Show more details of test coverage
* `make test` Run project test suite

### Creating release builds

Travis should build Linux and MacOS builds and upload to github releases

If you do it manually, you must follow the next steps:

* `make prepare-release` Change the resources in the app and prepare to release the app
* `make clean` Remove temporary files
* `make build` Build FiberCrypto Wallet
* Compress the content in `deploy` folder and inside that folder 


## WIP
This is a Work-In-Progress.
