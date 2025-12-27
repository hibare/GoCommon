# DB Package Documentation

## Overview

The `db` package provides database abstraction and utilities, supporting both PostgreSQL and SQLite. It offers interfaces for database operations, migration utilities, and test helpers for working with databases in Go applications.

---

## Key Types and Interfaces

- **Database**: Interface defining methods for database operations (`Open`, `Migrate`).
- **DB**: Struct wrapping a `gorm.DB` instance and its configuration.
- **DatabaseConfig**: Struct holding configuration for database connections (DSN, migrations path, etc.).
- **PostgresDatabase / SQLiteDatabase**: Implementations of the `Database` interface for PostgreSQL and SQLite.

---

## Main Functions

- **NewClient(ctx, config) (\*DB, error)**: Returns a singleton instance of the database connection.
- **(\*DB) Migrate() error**: Runs database migrations.
- **(\*DB) Close() error**: Closes the database connection.
- **(\*DB) RunSQLFromDirectory(dir string) error**: Executes all `.sql` files found in the specified directory in alphabetical order.
- **(\*DB) RunSQLFromFS(fsys fs.FS, dir string) error**: Executes all `.sql` files from an embedded filesystem directory in alphabetical order.
- **SetupMockPostgresDB()**: Sets up a mock PostgreSQL database for testing.
- **UnsetMockPostgresDB(container)**: Tears down the mock PostgreSQL database.

---

## Example Usage

### Basic Database Setup

```go
import (
    "context"
    "github.com/hibare/GoCommon/v2/pkg/db"
)

config := db.DatabaseConfig{
    DSN: "file::memory:?cache=shared",
    DBType: &db.SQLiteDatabase{},
}
client, err := db.NewClient(context.Background(), config)
if err != nil {
    panic(err)
}
err = client.Migrate()
```

### Running SQL Scripts from Directory

```go
// Execute all .sql files from a directory (e.g., for views, initial data, etc.)
err = client.RunSQLFromDirectory("path/to/sql/scripts")
if err != nil {
    panic(err)
}

// Or from embedded filesystem
//go:embed sql/scripts/*.sql
var sqlScripts embed.FS

err = client.RunSQLFromFS(sqlScripts, "sql/scripts")
```

---

## Notes

- Uses [gorm](https://gorm.io/) for ORM and [golang-migrate](https://github.com/golang-migrate/migrate) for migrations.
- Supports test helpers for both SQLite and PostgreSQL.
- Designed for extensibility and testability.
