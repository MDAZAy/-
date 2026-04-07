#!/usr/bin/env bash
set -euo pipefail

SOURCE="/opt/vpn-bot/shared/data/app.db"
TARGET_DIR="/opt/vpn-bot/backups"

mkdir -p "$TARGET_DIR"
test -f "$SOURCE"
cp "$SOURCE" "$TARGET_DIR/app-$(date +%F-%H%M%S).db"
find "$TARGET_DIR" -type f -name 'app-*.db' -mtime +14 -delete

