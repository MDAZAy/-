#!/usr/bin/env bash
set -euo pipefail

python -m venv .venv
. .venv/bin/activate
pip install -r embeddings/requirements.txt
uvicorn embeddings.embed_server:app --host 0.0.0.0 --port 8081
