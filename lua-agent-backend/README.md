# Lua Agent Backend

Backend-часть локального AI-агента для генерации Lua-кода без внешних LLM API.

## Соответствие ТЗ

- Локальная open-source модель через Ollama
- Без внешних AI-вендоров в runtime
- Генерация Lua-кода по естественному языку
- Проверка результата через validator pipeline
- Одна итерация self-correction при неуспешной валидации
- История генераций и статистика в PostgreSQL

## Рекомендуемая модель

Для демо-контуров под ограничение `8 GB VRAM`:

```bash
ollama pull qwen2.5-coder:7b
ollama create mws-agent -f ollama/Modelfile
```

Параметры запуска зафиксированы в `config/config.yaml` и соответствуют ограничению кейса:

- `num_ctx=4096`
- `num_predict=256`
- `batch=1`
- `parallel=1`

## Структура

- `cmd/agent` — HTTP сервер и graceful shutdown
- `config/config.yaml` — конфиг сервиса
- `internal/api` — `POST /generate`, `GET /history`, `GET /stats`, `GET /health`
- `internal/agent` — plan -> generate -> validate -> correct
- `internal/llm` — клиенты к Ollama и embed API
- `internal/storage` — PostgreSQL storage layer
- `internal/validator` — syntax, security, MWS, sandbox, pipeline
- `embeddings` — локальный Python embed server
- `ollama` — Modelfile для локальной модели
- `prompts` — текстовые промпты и few-shot набор
- `web` — простой demo UI
- `scripts` — запуск, benchmark и demo-запросы
- `metrics` — шаблоны отчётов для жюри
- `test/test_set.json` — набор задач для локальной проверки
- `migrations` — SQL миграции
- `docs/lua-agent-local-llm-task.md` — исходное ТЗ

## Быстрый старт без Docker

```powershell
cd lua-agent-backend
go mod tidy
go test ./...
go run ./cmd/agent
```

Перед запуском backend должны быть доступны:

- PostgreSQL с `pgvector`
- Ollama на хосте
- embed server из `embeddings/embed_server.py`

## Быстрый старт через Docker Compose

```bash
cp .env.example .env
docker compose up -d postgres embed
ollama pull qwen2.5-coder:7b
ollama create mws-agent -f ollama/Modelfile
docker compose up --build backend
```

В compose:

- `postgres` поднимается внутри контейнера
- `embed` поднимается внутри контейнера
- `backend` поднимается внутри контейнера
- `ollama` ожидается на хосте и доступен контейнеру через `host.docker.internal:11434`

## Demo UI

После старта backend открой:

- `http://127.0.0.1:8080/` — web demo
- `http://127.0.0.1:8080/health` — health-check компонентов
- `http://127.0.0.1:8080/history` — последние успешные генерации
- `http://127.0.0.1:8080/stats` — агрегированная статистика

## API

### `POST /generate`

```json
{
  "session_id": "optional-session-id",
  "prompt": "Напиши Lua-скрипт, который суммирует массив чисел"
}
```

### `GET /history?limit=10`

Возвращает последние успешные генерации.

### `GET /stats`

Возвращает агрегированную статистику по генерациям.

## Локальная валидация

В pipeline уже есть:

- синтаксическая проверка Lua через `gopher-lua`
- блокировка `os.execute`, `io.open`, `dofile`, `loadfile`
- MWS-проверки на `wf.vars`, `wf.initVariables`, `_utils.array`
- sandbox-выполнение с timeout и capture `print`
- одна итерация self-correction при неуспешной валидации

## Скрипты

- `scripts/run-local.sh` — локальный запуск backend
- `scripts/run-local.ps1` — локальный запуск backend на Windows
- `scripts/run_embed.sh` — запуск embed server
- `scripts/demo-request.ps1` — пример запроса к API
- `scripts/benchmark.sh` — простой benchmark-контур

## Что ещё можно улучшать

- реальный sentence-transformers embedder вместо hashed fallback
- более сильный benchmark harness по всему `test/test_set.json`
- загрузка `few_shot.json` в PostgreSQL как стартовых примеров
- более строгие интеграционные тесты с Postgres и Ollama
