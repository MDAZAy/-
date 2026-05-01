from __future__ import annotations

import json
import sys
import urllib.request
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]
TEST_SET_PATH = REPO_ROOT / "test" / "test_set.json"
BACKEND_URL = "http://127.0.0.1:8080/generate"


def call_backend(prompt: str) -> dict:
    payload = json.dumps({"prompt": prompt}).encode("utf-8")
    request = urllib.request.Request(
        BACKEND_URL,
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(request, timeout=120) as response:
        return json.loads(response.read().decode("utf-8"))


def main() -> None:
    test_set = json.loads(TEST_SET_PATH.read_text(encoding="utf-8"))
    results = []
    success = 0

    for item in test_set:
        response = call_backend(item["prompt"])
        ok = bool(response.get("validation", {}).get("ok"))
        if ok:
            success += 1
        results.append(
            {
                "id": item["id"],
                "prompt": item["prompt"],
                "ok": ok,
                "needs_clarification": bool(response.get("needs_clarification")),
                "model": response.get("model"),
            }
        )
        print(f"[{item['id']}] ok={ok} clarification={response.get('needs_clarification', False)}")

    summary = {
        "total": len(test_set),
        "success": success,
        "failed": len(test_set) - success,
        "success_rate": success / len(test_set) if test_set else 0,
        "results": results,
    }

    print(json.dumps(summary, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"run_test_set failed: {exc}", file=sys.stderr)
        raise
