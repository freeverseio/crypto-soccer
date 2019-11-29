CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('verse', '0');

CREATE TABLE tactics (
    team_id TEXT NOT NULL,
    verse BIGINT NOT NULL,

    tactic_id INT NOT NULL,

    shirt_0 INT CHECK (shirt_0 >= 0),
    shirt_1 INT CHECK (shirt_1 >= 0),
    shirt_2 INT CHECK (shirt_2 >= 0),
    shirt_3 INT CHECK (shirt_3 >= 0),
    shirt_4 INT CHECK (shirt_4 >= 0),
    shirt_5 INT CHECK (shirt_5 >= 0),
    shirt_6 INT CHECK (shirt_6 >= 0),
    shirt_7 INT CHECK (shirt_7 >= 0),
    shirt_8 INT CHECK (shirt_8 >= 0),
    shirt_9 INT CHECK (shirt_9 >= 0),
    shirt_10 INT CHECK (shirt_10 >= 0),
    shirt_11 INT CHECK (shirt_11 >= 0),
    shirt_12 INT CHECK (shirt_12 >= 0),
    shirt_13 INT CHECK (shirt_13 >= 0),

    extra_attack_1  BOOLEAN DEFAULT FALSE,
    extra_attack_2  BOOLEAN DEFAULT FALSE,
    extra_attack_3  BOOLEAN DEFAULT FALSE,
    extra_attack_4  BOOLEAN DEFAULT FALSE,
    extra_attack_5  BOOLEAN DEFAULT FALSE,
    extra_attack_6  BOOLEAN DEFAULT FALSE,
    extra_attack_7  BOOLEAN DEFAULT FALSE,
    extra_attack_8  BOOLEAN DEFAULT FALSE,
    extra_attack_9  BOOLEAN DEFAULT FALSE,
    extra_attack_10 BOOLEAN DEFAULT FALSE,

    PRIMARY KEY(team_id, verse)
)
