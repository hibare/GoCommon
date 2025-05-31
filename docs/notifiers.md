# Notifiers Package Documentation

## Overview

The `notifiers` package provides utilities for sending notifications to external services. The main subpackage is `discord`, which allows sending messages to Discord webhooks with rich embeds and components.

---

## Subpackages

- **discord**: Utilities for sending messages to Discord webhooks, including support for embeds, footers, and custom HTTP clients.

---

## Key Types and Functions (discord)

- **Client**: Interface for sending messages to Discord webhooks.
- **Options**: Struct for configuring the Discord client (webhook URL, HTTP client).
- **Message, Embed, EmbedField, ...**: Types for building Discord messages and embeds.
- **NewClient(opts) (Client, error)**: Creates a new Discord client.
- **(\*Client) Send(ctx, msg) error**: Sends a message to the Discord webhook.
- **AddFooter(footerStr) error**: Adds a footer to the last embed in a message.

---

## Example Usage

```go
import (
    "context"
    "github.com/hibare/GoCommon/v2/pkg/notifiers/discord"
)

client, _ := discord.NewClient(discord.Options{WebhookURL: "https://discord.com/api/webhooks/..."})
msg := &discord.Message{Content: "Hello, Discord!"}
_ = client.Send(context.Background(), msg)
```

---

## Notes

- Supports dependency injection for HTTP clients (testable/mocking).
- Rich embed support for Discord notifications.
