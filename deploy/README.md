# VPS Deploy

Этот набор рассчитан на Ubuntu/Debian VPS и структуру:

- `/opt/vpn-bot/current` — текущий релиз
- `/opt/vpn-bot/shared/data/app.db` — SQLite база
- `vpnbot` — системный пользователь для сервисов

## Что есть

- `env/*.production.env.example` — production шаблоны env
- `systemd/*.service` — сервисы backend и bot
- `nginx/*.conf` — nginx конфиги
- `scripts/install-vps.sh` — базовая подготовка сервера
- `scripts/deploy.sh` — выкладка и рестарт
- `scripts/smoke.sh` — простая post-deploy проверка
- `scripts/backup-db.sh` — backup SQLite с ротацией

## Базовый порядок

1. Подготовить VPS:
   - `sudo bash deploy/scripts/install-vps.sh`
2. Скопировать проект на сервер.
3. Создать env:
   - `cp deploy/env/backend.production.env.example /opt/vpn-bot/current/backend-go/.env`
   - `cp deploy/env/bot.production.env.example /opt/vpn-bot/current/bot-python/.env`
4. Заполнить реальные значения:
   - домен
   - `BOT_TOKEN`
   - `ADMIN_TOKEN`
   - затем YooKassa/VPN provider при подключении real integrations
5. Установить systemd unit-файлы в `/etc/systemd/system/`
6. Установить nginx конфиг в `/etc/nginx/sites-available/`
7. Выполнить деплой:
   - `sudo bash deploy/scripts/deploy.sh`
8. Проверить:
   - `sudo bash deploy/scripts/smoke.sh`

## Минимальные команды проверки

```bash
systemctl status vpn-backend.service
systemctl status vpn-bot.service
journalctl -u vpn-backend.service -n 100 --no-pager
journalctl -u vpn-bot.service -n 100 --no-pager
curl -fsS http://127.0.0.1:8080/health
```

## Важно

- SQLite база вынесена в `/opt/vpn-bot/shared/data/app.db`, чтобы переживать перевыкладки
- backend запускается из собранного бинарника, а не через `go run`
- для HTTPS используй `deploy/nginx/vpn-bot.ssl.conf` и certbot
