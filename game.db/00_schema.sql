CREATE TABLE player_props (
    player_id TEXT NOT NULL,
    player_name TEXT NOT NULL,
    PRIMARY KEY(player_id)
);

CREATE TABLE team_props (
    team_id TEXT NOT NULL,
    team_name TEXT NOT NULL DEFAULT '',
    team_manager_name TEXT NOT NULL DEFAULT '',
    maximum_bid BIGINT NOT NULL DEFAULT 10,
    PRIMARY KEY(team_id)
);
