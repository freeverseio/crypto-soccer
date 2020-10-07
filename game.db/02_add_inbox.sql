CREATE TYPE inbox_category AS ENUM ('offer', 'auction', 'promo', 'news', 'incident', 'welcome');

CREATE TABLE inbox (
    id BIGSERIAL,
    destinatary TEXT NOT NULL,
    category inbox_category NOT NULL,
    auction_id TEXT,
    title TEXT NOT NULL,
    text_message TEXT NOT NULL,
    custom_image_url TEXT,
    metadata JSON,
    is_read boolean DEFAULT false, 
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

CREATE INDEX idx_inbox_destinatary ON inbox (destinatary);
CREATE INDEX idx_inbox_created_at ON inbox (created_at);
