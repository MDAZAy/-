CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    last_prompt TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS histories (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_prompt TEXT NOT NULL,
    clarified_prompt TEXT NOT NULL DEFAULT '',
    generated_code TEXT NOT NULL,
    validation_status TEXT NOT NULL DEFAULT 'unknown',
    validation_errors JSONB NOT NULL DEFAULT '[]'::jsonb,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    model_name TEXT NOT NULL DEFAULT '',
    latency_ms BIGINT NOT NULL DEFAULT 0,
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_histories_session_id ON histories(session_id);
CREATE INDEX IF NOT EXISTS idx_histories_success_created_at ON histories(success, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_histories_validation_status ON histories(validation_status);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_sessions_updated_at ON sessions;
CREATE TRIGGER trg_sessions_updated_at
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_histories_updated_at ON histories;
CREATE TRIGGER trg_histories_updated_at
BEFORE UPDATE ON histories
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
