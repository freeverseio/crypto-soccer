CREATE TABLE player_props (
    player_id TEXT NOT NULL,
    player_name TEXT NOT NULL,
    PRIMARY KEY(player_id)
);

CREATE TABLE team_props (
    team_id TEXT NOT NULL,
    team_name TEXT NOT NULL,
    team_manager_name TEXT NOT NULL,
    PRIMARY KEY(team_id)
);

CREATE TABLE league_props (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    leaderboard JSONB,
    PRIMARY KEY(timezone_idx,country_idx, league_idx)
);