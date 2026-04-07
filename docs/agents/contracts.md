# Contracts And Rules

## HTTP API

Обязательные endpoints MVP:

- `GET /health`
- `POST /api/v1/users/ensure`
- `GET /api/v1/plans`
- `POST /api/v1/plans`
- `POST /api/v1/subscriptions/create`
- `GET /api/v1/subscriptions/active/:user_id`
- `POST /api/v1/payments/create`
- `POST /api/v1/payments/webhook`
- `POST /api/v1/vpn/issue`

## Правила для bot-python

- Только HTTP к backend.
- Никаких прямых SQL/ORM зависимостей.
- Никаких секретов платежей или VPN-провайдера.
- Если появляется новая кнопка, она должна опираться на существующий или новый backend endpoint.

## Правила для backend-go

- Новая бизнес-логика идёт в `services`.
- Хендлер не должен дублировать правила предметной области.
- Доступ к данным идёт через `repositories`.
- Реальные платежные/VPN интеграции должны подменять mock через provider interface.

## Security baseline

- Admin routes только через токен.
- `.env` не коммитится.
- Webhook принимает backend, не бот.

