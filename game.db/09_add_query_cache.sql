CREATE TABLE query_cache (
    key TEXT NOT NULL,
    data jsonb,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(key)
);

ALTER TYPE inbox_category ADD VALUE 'ban' AFTER 'welcome';
ALTER TYPE inbox_category ADD VALUE 'permaban' AFTER 'welcome';
