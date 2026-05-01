from __future__ import annotations

import json
import os
from pathlib import Path

import psycopg

from embed_server import hashed_embedding


REPO_ROOT = Path(__file__).resolve().parents[1]
FEW_SHOT_PATH = REPO_ROOT / "prompts" / "few_shot.json"
DATABASE_URL = os.getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/lua_agent?sslmode=disable")


def main() -> None:
    examples = json.loads(FEW_SHOT_PATH.read_text(encoding="utf-8"))
    with psycopg.connect(DATABASE_URL) as conn:
        with conn.cursor() as cur:
            for item in examples:
                task = item["task"]
                code = item["code"]
                vector = hashed_embedding(task).tolist()
                vector_literal = "[" + ",".join(f"{value:.8f}" for value in vector) + "]"
                cur.execute(
                    """
                    INSERT INTO histories (
                        id,
                        session_id,
                        user_prompt,
                        clarified_prompt,
                        generated_code,
                        validation_status,
                        validation_errors,
                        success,
                        model_name,
                        latency_ms,
                        input_tokens,
                        output_tokens,
                        metadata,
                        embedding
                    )
                    VALUES (
                        gen_random_uuid(),
                        gen_random_uuid(),
                        %s,
                        %s,
                        %s,
                        'passed',
                        '[]'::jsonb,
                        TRUE,
                        'few-shot-seed',
                        0,
                        0,
                        0,
                        '{}'::jsonb,
                        %s
                    )
                    ON CONFLICT DO NOTHING
                    """,
                    (task, task, code, vector_literal),
                )
        conn.commit()


if __name__ == "__main__":
    main()
