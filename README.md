<div align="center">
  <img src="./assets/logo.png" alt="GoCommon Logo" width="200" height="200">

[![Go Report Card](https://goreportcard.com/badge/github.com/hibare/GoCommon)](https://goreportcard.com/report/github.com/hibare/GoCommon)
[![GitHub issues](https://img.shields.io/github/issues/hibare/GoCommon)](https://github.com/hibare/GoCommon/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/hibare/GoCommon)](https://github.com/hibare/GoCommon/pulls)
[![GitHub](https://img.shields.io/github/license/hibare/GoCommon)](https://github.com/hibare/GoCommon/blob/main/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/hibare/GoCommon)](https://github.com/hibare/GoCommon/releases)

</div>

A modern, modular Go library providing reusable utilities and abstractions for cloud, configuration, concurrency, file operations, HTTP, logging, validation, and more. Designed for extensibility, testability, and ease of use in real-world Go projects.

---

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Makefile Commands](#makefile-commands)
- [Usage](#usage)
- [Package Documentation](#package-documentation)
  - [AWS](docs/aws.md)
  - [Concurrency](docs/concurrency.md)
  - [Config](docs/config.md)
  - [Constants](docs/constants.md)
  - [Context](docs/context.md)
  - [Crypto](docs/crypto.md)
  - [Datetime](docs/datetime.md)
  - [DB](docs/db.md)
  - [Env](docs/env.md)
  - [Errors](docs/errors.md)
  - [File](docs/file.md)
  - [HTTP](docs/http.md)
  - [Less](docs/less.md)
  - [Logger](docs/logger.md)
  - [Maps](docs/maps.md)
  - [Net](docs/net.md)
  - [Notifiers](docs/notifiers.md)
  - [Slice](docs/slice.md)
  - [Structs](docs/structs.md)
  - [Testhelper](docs/testhelper.md)
  - [Utils](docs/utils.md)
  - [Validator](docs/validator.md)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

GoCommon is a collection of well-tested, production-ready Go packages for common application needs:

- **Cloud**: AWS S3 utilities
- **Configuration**: YAML config loading, OS-aware paths
- **Concurrency**: Parallel task execution
- **Database**: Abstractions for PostgreSQL and SQLite
- **Environment**: Type-safe env var access
- **Errors**: Standardized error types
- **File**: Archiving, hashing, downloading, and more
- **HTTP**: Middleware, handlers, and client utilities
- **Logging**: Structured logging with multiple modes
- **Maps/Slices/Structs**: Generic helpers for data manipulation
- **Validation**: Struct validation with custom error messages
- **Versioning**: GitHub release checks

Each package is documented and designed for easy integration and testing.

---

## Installation

```sh
go get github.com/hibare/GoCommon/v2
```

---

## Makefile Commands

This project provides a `Makefile` to simplify common development tasks. Use these commands to lint, test, and manage dependencies:

| Command                      | Description                                  |
| ---------------------------- | -------------------------------------------- |
| `make help`                  | Display all available make targets           |
| `make clean`                 | Cleanup and tidy Go modules                  |
| `make init`                  | Initialize project (lint & pre-commit setup) |
| `make install-golangci-lint` | Install `golangci-lint` if not present       |
| `make install-pre-commit`    | Install pre-commit git hooks                 |
| `make test`                  | Run tests                                    |

**Example:**

```sh
make test
```

---

## Usage

Import only the packages you need. Example:

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/logger"
    "github.com/hibare/GoCommon/v2/pkg/env"
)

logger.InitDefaultLogger()
port := env.MustInt("PORT", 8080)
```

See the [Package Documentation](#package-documentation) for detailed usage and examples for each module.

---

## Package Documentation

Each package has its own detailed documentation:

- [AWS](docs/aws.md)
- [Concurrency](docs/concurrency.md)
- [Config](docs/config.md)
- [Constants](docs/constants.md)
- [Context](docs/context.md)
- [Crypto](docs/crypto.md)
- [Datetime](docs/datetime.md)
- [DB](docs/db.md)
- [Env](docs/env.md)
- [Errors](docs/errors.md)
- [File](docs/file.md)
- [HTTP](docs/http.md)
- [Less](docs/less.md)
- [Logger](docs/logger.md)
- [Maps](docs/maps.md)
- [Net](docs/net.md)
- [Notifiers](docs/notifiers.md)
- [Slice](docs/slice.md)
- [Structs](docs/structs.md)
- [Testhelper](docs/testhelper.md)
- [Utils](docs/utils.md)
- [Validator](docs/validator.md)

---

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, new features, or documentation improvements.

- Follow Go best practices and style guidelines.
- Add or update tests for your changes.
- Document new features in the relevant package doc.

---

## License

This project is licensed under the [MIT License](LICENSE).
