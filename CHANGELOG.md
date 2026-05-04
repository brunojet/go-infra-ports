# Changelog

## [1.1.0](https://github.com/brunojet/go-infra-ports/compare/v1.0.2...v1.1.0) (2026-05-04)


### Features

* **rest:** implement RestRepository core, registry, helpers and standardize tests ([98f0a54](https://github.com/brunojet/go-infra-ports/commit/98f0a542442c3fa07c881219d0a9b0e1629c03f7))
* **services:** add REST service layer with public exports ([d9f62ca](https://github.com/brunojet/go-infra-ports/commit/d9f62ca8a35c6b2fe78f951a470b809c18af00bb))
* **services:** add REST service layer with public exports ([8f44021](https://github.com/brunojet/go-infra-ports/commit/8f4402182ed8d32e1dc078231c0f77fd708a9f46))


### Bug Fixes

* **rest-service:** guard nil ctx pointers and init query/headers ([e7fb3b8](https://github.com/brunojet/go-infra-ports/commit/e7fb3b83e74ca28b8c0e2d78acd7b918e91ad87d))


### Documentation

* **services:** align Save semantics with full PUT replacement ([a26ee0d](https://github.com/brunojet/go-infra-ports/commit/a26ee0d282daa146428e1f5f7a3fffb5e12bd9a9))


### Code Refactoring

* **rest:** sanitize API — move RestMethod constants to public contracts, accept typed RestMethod in public API, remove dead code, update tests ([97f41f2](https://github.com/brunojet/go-infra-ports/commit/97f41f2ba46d2ad7ed2b22c97b4bc984c1c7ab36))

## [1.0.2](https://github.com/brunojet/go-infra-ports/compare/v1.0.1...v1.0.2) (2026-05-02)


### Bug Fixes

* **types:** remove obsolete http type file ([9691c20](https://github.com/brunojet/go-infra-ports/commit/9691c206e3a0c38a126b04c2232a0193d019b29f))

## [1.0.1](https://github.com/brunojet/go-infra-ports/compare/v1.0.0...v1.0.1) (2026-05-02)


### Bug Fixes

* **types:** remove obsolete http type file ([2385f85](https://github.com/brunojet/go-infra-ports/commit/2385f8568639e5a74ab532476c07d9a4807afdfb))
* **types:** remove obsolete http type file ([812531a](https://github.com/brunojet/go-infra-ports/commit/812531a3982617b766aa96a895621a170a814ebb))

## 1.0.0 (2026-05-02)


### Bug Fixes

* align coverage.sh with CI and handle empty package list in pre-commit ([b26f6cc](https://github.com/brunojet/go-infra-ports/commit/b26f6cc910b6e751d812ff8db2d384d9a928d79c))
