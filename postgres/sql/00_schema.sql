CREATE TABLE teams (
    id INT,
    name TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE players (
    id INT,
    name TEXT,
    team INT REFERENCES teams(id),
    PRIMARY KEY(id)
);