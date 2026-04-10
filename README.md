# VPN Telegram SaaS Bot

MVP-реализация проекта из `VPN_TZ_DD.txt`:

- `backend-go` — Go backend с SQLite, HTTP API, mock/CloudPayments payments, mock/3x-ui VPN и админкой
- `bot-python` — Telegram bot на `aiogram`, работающий только через HTTP API backend
- `deploy` — базовые артефакты для `systemd`, `nginx` и бэкапов
- `docs/agents` — контекст и процесс для агентской разработки

## Быстрый старт

### 1. Backend

```bash
cd backend-go
cp .env.example .env
go mod tidy
go run ./cmd/server
```

Backend поднимется на `http://localhost:8080`.

Локальные команды для Windows / PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/run-backend.ps1
powershell -ExecutionPolicy Bypass -File scripts/smoke-backend.ps1
powershell -ExecutionPolicy Bypass -File scripts/stop-backend.ps1
```

Полезные URL:

- `GET /health`
- `GET /admin?token=change-me-admin-token`

### 2. Bot

```bash
cd bot-python
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
cp .env.example .env
python -m app.main
```

Или через PowerShell-скрипт:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/check-bot.ps1
powershell -ExecutionPolicy Bypass -File scripts/run-bot.ps1
```

## MVP-покрытие

- `/start`, `/help`, `/menu`
- регистрация пользователя через `POST /api/v1/users/ensure`
- просмотр тарифов
- создание mock-платежа или CloudPayments payment link
- mock webhook / CloudPayments webhook с автосозданием подписки
- выдача mock VPN-ключа или 3x-ui VLESS Reality ссылки
- админ-панель со списками пользователей, тарифов, подписок, платежей и ключей
- job для истечения подписок и деактивации ключей

## Что пока заглушено

- реальный 3x-ui
- уведомления о продлении
- полный CRUD админки

## VPS

Для выкладки на сервер смотри:

- [deploy/README.md](c:\Projecs\VPN_BOT\deploy\README.md)
- [docs/deploy-vps.md](c:\Projecs\VPN_BOT\docs\deploy-vps.md)
