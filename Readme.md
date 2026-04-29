# go-logger

Minimal wrapper around `zap` with opinionated defaults: JSON, RFC3339 timestamps, stdout output, `caller`, and stack traces starting at `error` level. Common fields: `service`, `env`, and a required `module` when creating a logger.

## Install
```bash
go get github.com/NlightN22/go-logger@v1
```

## Example

``` go
// main.go
package main

import (
	"github.com/NlightN22/go-logger"
)

func main() {
	_ = logger.Init("info", "billing-api", "prod")
	log := logger.New("http")
	log.Infow("server started", "port", 8080)
	defer logger.Sync()
}
```
