# Suggested Agent Roles

## Backend Agent

- Работает только в `backend-go`
- Отвечает за модели, репозитории, сервисы, хендлеры, маршруты
- Любой новый процесс оплаты/VPN начинает с provider abstraction

## Bot Agent

- Работает только в `bot-python`
- Отвечает за Telegram UX, команды, кнопки, тексты
- Все новые сценарии заводит через `app/client.py`

## Deploy Agent

- Работает только в `deploy`
- Не меняет предметную логику
- Отвечает за сервисы, reverse proxy, backup, env-подготовку

## Docs Agent

- Работает только в `docs/agents` и `README.md`
- Сверяет документацию с реальным кодом, а не наоборот

## Reviewer Agent

- Проверяет, что:
  - bot не полез в БД
  - backend не утащил Telegram UX в себя
  - mock и real provider-слои не перемешаны
