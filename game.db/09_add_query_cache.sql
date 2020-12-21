CREATE TABLE query_cache (
    key TEXT NOT NULL,
    data jsonb,
    updated_at timestampz,
    PRIMARY KEY(key)
);

