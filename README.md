# Seq Logger

## Overview

`go-seqlogger` is a custom log handler for sending structured logs to Seq using Go's `slog` package.

## Installation

```bash
go get github.com/swczk/go-seqlogger
```

## Usage

### Basic Configuration

```go
import (
    "log/slog"
    "github.com/swczk/go-seqlogger"
)

// Create default configuration
config := seqlogger.DefaultConfig("http://seq-example:5341").
    WithAPIKey("your-api-key").
    WithLogLevel(slog.LevelInfo).
    WithSourceTracking().
    WithRequestIDKey("request-id")

// Create logger
logger := seqlogger.New(config)

// Log messages
logger.Info("User logged in", "user_id", 123)
logger.Error("Failed to process request", "error", err)
```

## Configuration Options

- `Endpoint`: Seq server URL
- `APIKey`: Authentication key for Seq
- `LogLevel`: Minimum log level to send
- `AddSource`: Include source code location
- `RequestIDKey`: Context key for request tracing

## Features

- Structured logging
- Seq CLEF format compatibility
- Configurable log levels
- Source code tracking
- Request ID injection

## License

