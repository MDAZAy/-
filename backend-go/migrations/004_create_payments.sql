CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    provider TEXT NOT NULL,
    external_payment_id TEXT NOT NULL UNIQUE,
    payment_url TEXT,
    raw_response TEXT,
    created_at DATETIME NOT NULL
);

