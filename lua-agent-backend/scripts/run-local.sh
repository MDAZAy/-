#!/usr/bin/env bash
set -euo pipefail

export CONFIG_PATH="${CONFIG_PATH:-config/config.yaml}"
go run ./cmd/agent
