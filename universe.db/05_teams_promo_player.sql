CREATE TABLE  promo_players (
    team_id TEXT NOT NULL REFERENCES teams(team_id),
    claimed BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (team_id)
);