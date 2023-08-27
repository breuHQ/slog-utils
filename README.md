# slog-utils: utilities to work with slog

![GitHub release (with filter)](https://img.shields.io/github/v/release/breuHQ/slog-utils)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/breuHQ/slog-utils)
[![License](https://img.shields.io/github/license/breuHQ/slog-utils)](./LICENSE)
![GitHub contributors](https://img.shields.io/github/contributors/breuHQ/slog-utils)


Utilities for [slog](https://pkg.go.dev/log/slog) for

- `calldepth`, add caller skip for helper functions e.g. in 3rd party libraries e.g. [temporal](go.temporal.io/sdk) et el.

## [calldepth](./calldepth/): add caller skip, similar to that provided by

- zap provides [`AddCallerSkip`](https://pkg.go.dev/go.uber.org/zap#AddCallerSkip)—this plus [`WithOptions`](https://pkg.go.dev/go.uber.org/zap#Logger.WithOptions) allows for a simple adjustment in helpers [as shown here](https://go.dev/play/p/lMB0CJ1_E-3)
- zerolog provides [`Caller`](https://pkg.go.dev/github.com/rs/zerolog#Event.Caller) and [`CallerSkipFrame`](https://pkg.go.dev/github.com/rs/zerolog#Event.CallerSkipFrame)
- go-kit/log provides [`Caller`](https://pkg.go.dev/github.com/go-kit/log#Caller) (and [there's some discussion of making `Caller` even more helpful to users](https://github.com/go-kit/log/issues/16#issuecomment-1236410861))
- glog provides [several `X-Depth` functions](https://pkg.go.dev/github.com/golang/glog) that allow the user to adjust the depth
- the stdlib log provides [`Output`](https://pkg.go.dev/log#Logger.Output)

> code taken from [https://github.com/golang/go/issues/59145#issuecomment-1481920720](https://github.com/golang/go/issues/59145#issuecomment-1481920720)

### 🚀 Install

```sh
go get go.breu.io/slog-utils
```

**Compatibility**: go >= 1.21

- ⚠️ Work in Progress
- ⚠️ Use this library carefully, log processing can be very costly (!)

### 💡 Usage

```go
package main

import (
 "log/slog"

 "go.temporal.io/sdk/client"

 "go.breu.io/slog-utils/calldepth"
)

func main() {
  adapter := calldepth.NewAdapter(
    calldepth.NewLogger(slog.NewStdLogger()),
    calldepth.WithCallDepth(5), // 5 for activities, 6 for workflows.
  )
  opts := client.Options{
    HostPort: "localhost:7233",
    Logger:   adapter.WithGroup("temporal"),
  }

  tclient, err := client.Dial(opts)
}
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=breuHQ/slog-utils)


## 📝 License

Copyright © 2023 [Breu Inc.](https://github.com/breuHQ)

This project is [MIT](./LICENSE) licensed.
