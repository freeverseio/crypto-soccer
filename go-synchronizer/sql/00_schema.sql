CREATE TABLE teams (
    id INT,
    name TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '-1');

CREATE TABLE players (
    id INT,
    name TEXT,
    team INT REFERENCES teams(id),
    PRIMARY KEY(id)
);
