#!/usr/bin/env bash
set -euo pipefail

cp -n .env.example .env || true
docker compose up -d postgres embed
echo "Install Ollama on the host and run:"
echo "  ollama pull qwen2.5-coder:7b"
echo "  ollama create mws-agent -f ollama/Modelfile"
