#!/usr/bin/env bash
set -euo pipefail

PROMPT="${1:-Write Lua code that sums an array of integers}"
START=$(date +%s)
curl -s http://127.0.0.1:8080/generate \
  -H "Content-Type: application/json" \
  -d "{\"prompt\":\"${PROMPT}\"}" > /tmp/lua-agent-benchmark.json
END=$(date +%s)
echo "Elapsed seconds: $((END - START))"
echo "Response saved to /tmp/lua-agent-benchmark.json"
echo "Measure VRAM separately with: watch -n 0.5 nvidia-smi"
