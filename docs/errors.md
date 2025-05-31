# Errors Package Documentation

## Overview

The `errors` package provides common error variables and a custom error type for the GoCommon library. It standardizes error handling and messaging across the codebase.

---

## Key Types and Variables

- **Error**: Struct representing an error with a code and message.
- **ErrUnauthorized**: Indicates an unauthorized action.
- **ErrInternalServerError**: Indicates an internal server error.
- **ErrNonOKError**: Indicates a non-OK error.
- **ErrChecksumMismatch**: Indicates a checksum mismatch.
- **ErrCreatingDir**: Indicates a directory creation error.
- **ErrNotDir**: Indicates a not-a-directory error.
- **ErrNotFile**: Indicates a not-a-file error.
- **ErrRecordNotFound**: Indicates a record not found error.

---

## Usage

Use these error variables for consistent error handling and messaging throughout your application.
