CREATE TABLE owner_props (
    owner TEXT NOT NULL,
    maximum_bid BIGINT NOT NULL DEFAULT 10,
    PRIMARY KEY(owner)
);

ALTER TABLE team_props ALTER COLUMN team_name DROP NOT NULL;
ALTER TABLE team_props ALTER COLUMN team_manager_name DROP NOT NULL;