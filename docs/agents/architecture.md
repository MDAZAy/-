# Architecture Map

## Контуры

### `backend-go`

- `cmd/server` — точка входа
- `internal/config` — env-конфиг
- `internal/db` — подключение, миграция, seed
- `internal/models` — сущности GORM
- `internal/dto` — контракт API
- `internal/repositories` — слой доступа к данным
- `internal/services` — бизнес-логика
- `internal/handlers` — HTTP-обработчики
- `internal/routes` — маршруты
- `internal/middleware` — логирование, recovery, admin auth
- `internal/web/templates` — admin panel и mock payment UI

### `bot-python`

- `app/config.py` — env-конфиг
- `app/client.py` — единственная точка HTTP-доступа к backend
- `app/handlers/*.py` — команды и меню Telegram
- `app/keyboards` — reply/inline клавиатуры

### `deploy`

- `nginx` — reverse proxy
- `systemd` — сервисы backend и bot
- `scripts` — выкладка и backup

## Инварианты

- Бот не трогает БД и не решает бизнес-логику.
- Оплата и VPN живут только в backend.
- Новые интеграции сначала добавляются через интерфейсы провайдеров в backend.
- Любой admin endpoint должен оставаться защищён токеном.

