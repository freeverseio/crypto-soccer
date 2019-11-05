CREATE TABLE tactics (
    team_id TEXT NOT NULL,
    verse BIGINT NOT NULL,

    tactic_id INT NOT NULL,

    shirt_0 INT,
    shirt_1 INT,
    shirt_2 INT,
    shirt_3 INT,
    shirt_4 INT,
    shirt_5 INT,
    shirt_6 INT,
    shirt_7 INT,
    shirt_8 INT,
    shirt_9 INT,
    shirt_10 INT,
    shirt_11 INT,
    shirt_12 INT,
    shirt_13 INT,

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
