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

// Works with After() as well - this example will attempt to wait for five
// seconds but will only wait one second since the mocked time is advanced ten
// seconds in the goroutine, triggering the channel send.
go func() {
    time.Sleep(1 * time.Second)
    fmt.Println("Advancing mocked time by 10 seconds")
    time.Advance(10 * time.Second)
}()
fmt.Println("Waiting 5 mocked seconds...")
<-mocktime.After(5 * time.Second)
fmt.Println("...time elapsed!")
```
