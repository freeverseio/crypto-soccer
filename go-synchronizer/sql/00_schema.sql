CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE timezones (
    id INT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE countries (
    id INT NOT NULL,
    timezone_id INT NOT NULL REFERENCES timezones(id),
    PRIMARY KEY(id)
);

CREATE TABLE leagues (
    id NUMERIC(78,0) NOT NULL,
    country_id INT NOT NULL REFERENCES countries(id),
    PRIMARY KEY(id)
);

CREATE TABLE teams (
    id NUMERIC(78,0) NOT NULL,
    league_id INT NOT NULL REFERENCES leagues(id),
    owner TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE players (
    id NUMERIC(78,0) NOT NULL,
    -- monthOfBirthInUnixTime TEXT NOT NULL,
    -- blockNumber BIGINT NOT NULL,
    -- inBlockIndex INT NOT NULL,
    team_id BIGINT NOT NULL REFERENCES teams(id),
    -- state TEXT NOT NULL,
    defence INT NOT NULL,
    speed INT NOT NULL,
    pass INT NOT NULL,
    shoot INT NOT NULL,
    endurance INT NOT NULL,
    PRIMARY KEY(id)
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



