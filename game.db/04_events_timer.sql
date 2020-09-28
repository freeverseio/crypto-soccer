CREATE TYPE entity_key_type AS ENUM ('offer', 'auction');

CREATE TABLE mailbox_cron (
    entity_key entity_key_type NOT NULL,
    last_time_checked timestamp without time zone NOT NULL,
);
