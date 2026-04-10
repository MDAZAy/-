# Lua Agent Backend

Отдельный backend-репозиторий для локального AI-агента, который генерирует Lua-код без внешних LLM API.

## Что внутри

- `internal/storage` — слой хранения истории генераций и статистики в PostgreSQL
- `migrations` — SQL-миграции для `sessions` и `histories`
- `docs/lua-agent-local-llm-task.md` — исходное ТЗ

## Быстрый старт

```powershell
cd lua-agent-backend
go mod tidy
go test ./...
```

## Текущий статус

Сейчас в репозитории реализован storage-слой:

- `Save`
- `GetRecentSuccess`
- `GetStats`

Следующие слои для разработки:

- `internal/validator`
- `internal/llm`
- `internal/agent`
- `internal/api`
- `cmd/agent`
