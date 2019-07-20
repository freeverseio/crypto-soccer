CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE teams (
    id BIGINT,
    name TEXT NOT NULL,
    creationTimestamp BIGINT NOT NULL,
    blockNumber BIGINT NOT NULL,
    owner TEXT NOT NULL,
    currentLeagueId BIGINT NOT NULL,
    posInCurrentLeagueId INT NOT NULL,
    prevLeagueId BIGINT NOT NULL,
    posInPrevLeagueId INT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE teams_history (
    teamId BIGINT NOT NULL REFERENCES teams(id),
    blockNumber BIGINT NOT NULL,
    owner TEXT NOT NULL,
    currentLeagueId BIGINT NOT NULL,
    posInCurrentLeagueId INT NOT NULL,
    prevLeagueId BIGINT NOT NULL,
    posInPrevLeagueId INT NOT NULL,
    PRIMARY KEY(teamId, blockNumber)
);

CREATE TABLE players (
    id BIGINT,
    monthOfBirthInUnixTime TEXT NOT NULL,
    blockNumber BIGINT NOT NULL,
    teamId BIGINT NOT NULL REFERENCES teams(id),
    state TEXT NOT NULL,
    defence INT NOT NULL,
    speed INT NOT NULL,
    pass INT NOT NULL,
    shoot INT NOT NULL,
    endurance INT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE players_history (
    playerId BIGINT NOT NULL REFERENCES players(id),
    blockNumber BIGINT NOT NULL,
    teamId BIGINT NOT NULL REFERENCES teams(id),
    state TEXT NOT NULL,
    defence INT NOT NULL,
    speed INT NOT NULL,
    pass INT NOT NULL,
    shoot INT NOT NULL,
    endurance INT NOT NULL,
    PRIMARY KEY(playerId, blockNumber)
);
