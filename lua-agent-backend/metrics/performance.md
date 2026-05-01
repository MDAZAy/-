# Performance

## Goal

Measure:

- generation latency
- backend response time
- peak VRAM during Ollama generation

## Commands

```bash
ollama ps
watch -n 0.5 nvidia-smi
./scripts/benchmark.sh "Write Lua code that sums an array of integers"
```

## Result Table

| Scenario | Model | Prompt | Time (s) | Peak VRAM (GB) | Notes |
| --- | --- | --- | --- | --- | --- |
| baseline | qwen2.5-coder:7b | array sum | TBD | TBD | fill during jury run |
