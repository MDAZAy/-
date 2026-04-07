#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://127.0.0.1:8080}"

echo "== health =="
curl -fsS "$BASE_URL/health"
echo
echo "== plans =="
curl -fsS "$BASE_URL/api/v1/plans"
echo
echo "== admin =="
curl -fsS -H "X-Admin-Token: ${ADMIN_TOKEN:-change-me-admin-token}" "$BASE_URL/admin" >/dev/null
echo "admin ok"

