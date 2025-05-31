# Net Package Documentation

## Overview

The `net` package provides utilities for network-related operations. The main subpackage is `ip`, which offers functions for working with IP addresses.

---

## Subpackages

- **ip**: Utilities for IP address operations, such as checking if an IP is public.

---

## Key Functions (ip)

- **IsPublicIP(ipStr string) bool**: Returns true if the given IP address is public (not loopback, link-local, or private).

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/net/ip"
)

isPublic := ip.IsPublicIP("8.8.8.8") // true
```

---

## Notes

- Useful for filtering or validating IP addresses in network applications.
