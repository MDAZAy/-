#!/usr/bin/env bash
set -euo pipefail

apt-get update
apt-get install -y nginx rsync curl python3 python3-venv ca-certificates

if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed. Install Go 1.22+ before running deploy.sh"
fi

systemctl enable nginx

