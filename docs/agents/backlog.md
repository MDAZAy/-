# Backlog After MVP

## Высокий приоритет

- Подключить реальный CloudPayments provider вместо mock.
- Подключить реальный x-ui/3x-ui provider вместо mock.
- Добавить историю платежей и подписок по пользователю.
- Добавить фильтры и действия в admin panel.

## Средний приоритет

- Напоминания о продлении.
- Уведомления об истечении подписки.
- Переход с SQLite на PostgreSQL.
- Наблюдаемость: structured logs, metrics, health dependencies.

## Технический долг

- Покрыть backend unit/integration тестами.
- Добавить CI для lint/build.
- Уточнить схему миграций и отказаться от `AutoMigrate` в production.
