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

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE players (
    id INT,
    state TEXT,
    defence INT NOT NULL DEFAULT 0,
    speed INT NOT NULL DEFAULT 0,
    pass INT NOT NULL DEFAULT 0,
    shoot INT NOT NULL DEFAULT 0,
    endurance INT NOT NULL DEFAULT 0,
    team INT REFERENCES teams(id),
    PRIMARY KEY(id)
);
