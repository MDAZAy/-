# Deploy To VPS

## Цель

Поднять `backend-go` и `bot-python` на Ubuntu/Debian VPS так, чтобы:

- backend был доступен локально на `127.0.0.1:8080`
- nginx принимал внешний трафик
- bot стабильно ходил в Telegram Bot API
- SQLite база жила вне текущего релиза

## Что подготовить заранее

- VPS с Ubuntu/Debian
- домен или поддомен, например `vpn.example.com`
- Go 1.22+
- Python 3.11+
- токен Telegram-бота
- `ADMIN_TOKEN`

## Файлы

- [deploy/README.md](/c:/Projecs/VPN_BOT/deploy/README.md)
- [deploy/env/backend.production.env.example](/c:/Projecs/VPN_BOT/deploy/env/backend.production.env.example)
- [deploy/env/bot.production.env.example](/c:/Projecs/VPN_BOT/deploy/env/bot.production.env.example)
- [deploy/systemd/vpn-backend.service](/c:/Projecs/VPN_BOT/deploy/systemd/vpn-backend.service)
- [deploy/systemd/vpn-bot.service](/c:/Projecs/VPN_BOT/deploy/systemd/vpn-bot.service)
- [deploy/nginx/vpn-bot.ssl.conf](/c:/Projecs/VPN_BOT/deploy/nginx/vpn-bot.ssl.conf)

## Короткий сценарий

1. Установить системные пакеты и nginx.
2. Скопировать проект в `/opt/vpn-bot/current`.
3. Заполнить `.env` для backend и bot.
   - для CloudPayments переключить `PAYMENT_PROVIDER=cloudpayments`
   - указать `CLOUDPAYMENTS_PUBLIC_ID` и `CLOUDPAYMENTS_API_SECRET`
   - в кабинете CloudPayments указать webhook URL: `https://vpn.example.com/api/v1/payments/webhook`
   - для 3x-ui переключить `VPN_PROVIDER=3xui`
   - указать `VPN_PROVIDER_ENDPOINT`, `VPN_PROVIDER_USERNAME`, `VPN_PROVIDER_PASSWORD`
   - указать `VPN_PROVIDER_INBOUND_ID`, `VPN_PROVIDER_PUBLIC_HOST`, `VPN_PROVIDER_PUBLIC_PORT`
   - указать `VPN_PROVIDER_REALITY_SERVER_NAME`, `VPN_PROVIDER_REALITY_PUBLIC_KEY`, `VPN_PROVIDER_REALITY_SHORT_ID`
4. Установить systemd unit-файлы.
5. Включить nginx конфиг.
6. Запустить `deploy/scripts/deploy.sh`.
7. Прогнать `deploy/scripts/smoke.sh`.

## После запуска

- backend health: `curl http://127.0.0.1:8080/health`
- admin panel: `https://vpn.example.com/admin?token=...`
- bot logs: `journalctl -u vpn-bot.service -f`

## Следующий шаг после VPS

После первого успешного VPS запуска уже имеет смысл:

- подключать real 3x-ui provider
- переводить bot с polling на более строгий production режим по твоему выбору
