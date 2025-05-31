# Constants Package Documentation

## Overview

The `constants` package provides application-wide constants for configuration, file naming, and formatting. These constants are used throughout the GoCommon library to ensure consistency and avoid magic strings.

---

## Main Constants

- **DefaultDateTimeLayout**: The default date-time layout string (e.g., `20060102150405`).
- **ConfigRootLinux**: Default config root directory for Linux (`/etc/`).
- **ConfigRootWindows**: Default config root directory for Windows (`C:\\ProgramData\\`).
- **ConfigRootDarwin**: Default config root directory for macOS (`/Library/Application Support/`).
- **ConfigFileName**: Default config file name (`config`).
- **ConfigFileExtension**: Default config file extension (`yaml`).
- **S3PrefixSeparator**: Default S3 prefix separator (`/`).

---

## Usage

These constants are imported and used by other packages for file paths, configuration, and AWS S3 operations.
