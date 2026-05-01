CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE histories
ADD COLUMN IF NOT EXISTS embedding vector(128);

CREATE INDEX IF NOT EXISTS idx_histories_embedding
ON histories
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
