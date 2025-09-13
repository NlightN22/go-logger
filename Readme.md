# go-logger

Минималистичная обёртка над `zap` с фиксированными стандартами: JSON, RFC3339, вывод в stdout, `caller`, стек с уровня `error`. Единые поля: `service`, `env`, и обязательный `module` при создании логгера.

## Установка
```bash
go get github.com/NlightN22/go-logger@v1
```

## Example

``` go
// main.go
package main

import (
	"github.com/NlightN22/go-logger@v1"
)

func main() {
	_ = logger.Init("info", "billing-api", "prod")
	log := logger.New("http")
	log.Infow("server started", "port", 8080)
	defer logger.Sync()
}
```