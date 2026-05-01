from __future__ import annotations

import hashlib
import math
import os
from typing import List

import numpy as np
from fastapi import FastAPI
from pydantic import BaseModel


EMBED_DIM = int(os.getenv("EMBED_DIM", "128"))

app = FastAPI(title="Lua Agent Embed Server", version="0.1.0")


class EmbedRequest(BaseModel):
    text: str


class EmbedResponse(BaseModel):
    vector: List[float]


def hashed_embedding(text: str, dim: int = EMBED_DIM) -> np.ndarray:
    vector = np.zeros(dim, dtype=np.float32)
    for token in text.lower().split():
        digest = hashlib.sha256(token.encode("utf-8")).digest()
        for index, byte in enumerate(digest):
            bucket = index % dim
            sign = -1.0 if byte % 2 else 1.0
            vector[bucket] += sign * (byte / 255.0)

    norm = math.sqrt(float(np.dot(vector, vector)))
    if norm > 0:
        vector /= norm
    return vector


@app.get("/health")
def health() -> dict:
    return {"status": "ok", "embed_dim": EMBED_DIM, "mode": "hashed-local"}


@app.post("/embed", response_model=EmbedResponse)
def embed(request: EmbedRequest) -> EmbedResponse:
    vector = hashed_embedding(request.text or "")
    return EmbedResponse(vector=vector.tolist())
