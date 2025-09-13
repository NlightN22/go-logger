# File: README.md

# go-logger

Минималистичная обёртка над `zap` с фиксированными стандартами: JSON, RFC3339, вывод в stdout, `caller`, стек с уровня `error`. Единые поля: `service`, `env`, и обязательный `module` при создании логгера.

## Установка
```bash
go get github.com/NlightN22/go-logger@v1
