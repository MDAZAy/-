# Jury Runbook

## 1. Model

```bash
ollama pull qwen2.5-coder:7b
ollama create mws-agent -f ollama/Modelfile
```

## 2. Start infrastructure

```bash
cp .env.example .env
docker compose up -d postgres embed
docker compose up --build backend
```

## 3. Verify health

```bash
curl http://127.0.0.1:8080/health
```

Expected:

- `storage.status = ok`
- `ollama.status = ok`

## 4. Run one demo request

```bash
curl -X POST http://127.0.0.1:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt":"Write Lua code that sums an array of integers and prints the result."}'
```

## 5. Run the local test set

```bash
python scripts/run_test_set.py
```

## 6. Measure VRAM

Run on the host while sending the request:

```bash
watch -n 0.5 nvidia-smi
```

Record the peak memory into:

- `metrics/performance.md`
- `metrics/comparison.md`
