CREATE TABLE IF NOT EXISTS vpn_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL,
    external_client_id TEXT NOT NULL,
    key_name TEXT NOT NULL,
    access_url TEXT NOT NULL,
    config_json TEXT,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    expires_at DATETIME,
    created_at DATETIME NOT NULL
);

