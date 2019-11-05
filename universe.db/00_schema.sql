CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE timezones (
    timezone_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx)
);

CREATE TABLE countries (
    timezone_idx INT NOT NULL REFERENCES timezones(timezone_idx),
    country_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx, country_idx)
);

CREATE TABLE leagues (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx,country_idx, league_idx),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx)
);

CREATE TABLE teams (
    team_id NUMERIC(78,0) NOT NULL,
    name TEXT NOT NULL,
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    owner TEXT NOT NULL,
    league_idx INT NOT NULL,
    team_idx_in_league INT NOT NULL,
    points INT NOT NULL DEFAULT 0,
    w INT NOT NULL DEFAULT 0,
    d INT NOT NULL DEFAULT 0,
    l INT NOT NULL DEFAULT 0,
    goals_forward INT NOT NULL DEFAULT 0,
    goals_against INT NOT NULL DEFAULT 0,
    PRIMARY KEY(team_id),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);

CREATE TABLE players (
    player_id NUMERIC(78,0) NOT NULL,
    name TEXT NOT NULL,
    team_id NUMERIC(78,0) NOT NULL REFERENCES teams(team_id),
    defence INT NOT NULL,
    speed INT NOT NULL,
    pass INT NOT NULL,
    shoot INT NOT NULL,
    endurance INT NOT NULL,
    shirt_number INT NOT NULL,
    preferred_position TEXT NOT NULL,
    potential INT NOT NULL, 
    --- TODO remove default
    age INT NOT NULL DEFAULT 0, 
    encoded_skills TEXT NOT NULL,
    encoded_state TEXT NOT NULL,
    frozen BOOLEAN NOT NULL,
    PRIMARY KEY(player_id)
);

CREATE TABLE matches (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    match_day_idx INT NOT NULL,
    match_idx INT NOT NULL,
    home_team_id NUMERIC(78,0) REFERENCES teams(team_id),
    visitor_team_id NUMERIC(78,0) REFERENCES teams(team_id),
    home_goals INT,
    visitor_goals INT,
    match_log_half_1 TEXT,
    match_log_half_2 TEXT,
    PRIMARY KEY(timezone_idx,country_idx, league_idx, match_day_idx, match_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);

-- CREATE TABLE teams_history (
--     teamId BIGINT NOT NULL REFERENCES teams(id),
--     blockNumber BIGINT NOT NULL,
--     inBlockIndex INT NOT NULL,
--     owner TEXT NOT NULL,
--     currentLeagueId BIGINT NOT NULL,
--     posInCurrentLeagueId INT NOT NULL,
--     prevLeagueId BIGINT NOT NULL,
--     posInPrevLeagueId INT NOT NULL,
--     PRIMARY KEY(teamId, blockNumber, inBlockIndex)
-- );

-- CREATE TABLE players (
--     id BIGINT NOT NULL,
--     monthOfBirthInUnixTime TEXT NOT NULL,
--     blockNumber BIGINT NOT NULL,
--     inBlockIndex INT NOT NULL,
--     teamId BIGINT NOT NULL REFERENCES teams(id),
--     state TEXT NOT NULL,
--     defence INT NOT NULL,
--     speed INT NOT NULL,
--     pass INT NOT NULL,
--     shoot INT NOT NULL,
--     endurance INT NOT NULL,
--     PRIMARY KEY(id)
-- );

-- CREATE TABLE players_history (
--     playerId BIGINT NOT NULL REFERENCES players(id),
--     blockNumber BIGINT NOT NULL,
--     inBlockIndex INT NOT NULL,
--     teamId BIGINT NOT NULL REFERENCES teams(id),
--     state TEXT NOT NULL,
--     defence INT NOT NULL,
--     speed INT NOT NULL,
--     pass INT NOT NULL,
--     shoot INT NOT NULL,
--     endurance INT NOT NULL,
--     PRIMARY KEY(playerId, blockNumber, inBlockIndex)
-- );



