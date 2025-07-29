## mocktime

[![Build Status](https://github.com/nitroshare/mocktime/actions/workflows/test.yml/badge.svg)](https://github.com/nitroshare/mocktime/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/nitroshare/mocktime/badge.svg?branch=main)](https://coveralls.io/github/nitroshare/mocktime?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/nitroshare/mocktime.svg)](https://pkg.go.dev/github.com/nitroshare/mocktime)
[![MIT License](https://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](https://opensource.org/licenses/MIT)

This package provides an easy way to mock specific functions in the `time` package:

```golang
import "github.com/nitroshare/mocktime"

// Same as time.Now()
mocktime.Now()

// Mock Now() and After()
mocktime.Mock()
defer mocktime.Unmock()

// All calls to Now() will return the same time...
mocktime.Now()

// ...until the time is advanced
mocktime.Advance(5 * time.Second)

// ...or explicitly set
mocktime.Set(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC))

// Calls to After() will block until the time is advanced (in another
// goroutine, for example)
<-mocktime.After(5 * time.Second)
```

### AdvanceToAfter

A special utility function is provided that blocks until the next call to `After()`:

```golang
chanDone := make(chan any)

go func() {

    // Normally, this will block until Set() or Advance() is called
    <-mocktime.After(5 * time.Second)

    close(chanDone)
}()

// Advance the time to the expiry of the After() call above; we don't need to
// worry if the goroutine has reached the After() call or not when this
// function is called as it will block until After() is called
AdvanceToAfter()

// This read is guaranteed to succeed because the read on After() in the
// goroutine is unblocked
<-chanDone
```
