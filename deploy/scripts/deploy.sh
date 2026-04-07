#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="/opt/vpn-bot"
APP_USER="vpnbot"
RELEASE_DIR="$ROOT_DIR/current"
SHARED_DIR="$ROOT_DIR/shared"
BACKEND_DIR="$RELEASE_DIR/backend-go"
BOT_DIR="$RELEASE_DIR/bot-python"

mkdir -p "$ROOT_DIR" "$SHARED_DIR/data" "$BACKEND_DIR/bin"
id -u "$APP_USER" >/dev/null 2>&1 || useradd --system --create-home --shell /usr/sbin/nologin "$APP_USER"

rsync -av --delete ./ "$RELEASE_DIR"/ \
  --exclude ".git" \
  --exclude ".venv" \
  --exclude "__pycache__" \
  --exclude ".cache"

chown -R "$APP_USER":"$APP_USER" "$ROOT_DIR"

cd "$BACKEND_DIR"
mkdir -p bin
go mod tidy
go build -o bin/backend ./cmd/server

cd "$BOT_DIR"
python3 -m venv .venv
. .venv/bin/activate
pip install --upgrade pip
pip install -r requirements.txt

test -f "$BACKEND_DIR/.env" || cp "$RELEASE_DIR/deploy/env/backend.production.env.example" "$BACKEND_DIR/.env"
test -f "$BOT_DIR/.env" || cp "$RELEASE_DIR/deploy/env/bot.production.env.example" "$BOT_DIR/.env"

systemctl daemon-reload
systemctl restart vpn-backend.service
systemctl restart vpn-bot.service
systemctl --no-pager --full status vpn-backend.service
systemctl --no-pager --full status vpn-bot.service

