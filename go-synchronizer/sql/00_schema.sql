CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE teams (
    id INT,
    name TEXT NOT NULL,
    creationTimestamp TEXT NOT NULL,
    blockNumber TEXT NOT NULL,
    owner TEXT NOT NULL,
    currentLeagueId INT NOT NULL,
    posInCurrentLeagueId INT NOT NULL,
    prevLeagueId INT NOT NULL,
    posInPrevLeagueId INT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE teams_history (
    teamId INT NOT NULL REFERENCES teams(id),
    blockNumber TEXT NOT NULL,
    owner TEXT NOT NULL,
    currentLeagueId INT NOT NULL,
    posInCurrentLeagueId INT NOT NULL,
    prevLeagueId INT NOT NULL,
    posInPrevLeagueId INT NOT NULL,
    PRIMARY KEY(teamId, blockNumber)
);

CREATE TABLE players (
    id INT,
    monthOfBirthInUnixTime TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE players_history (
    playerId INT NOT NULL REFERENCES players(id),
    blockNumber TEXT NOT NULL,
    teamId INT NOT NULL REFERENCES teams(id),
    state TEXT NOT NULL,
    defence INT NOT NULL,
    speed INT NOT NULL,
    pass INT NOT NULL,
    shoot INT NOT NULL,
    endurance INT NOT NULL,
    PRIMARY KEY(playerId, blockNumber)
);
