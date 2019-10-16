CREATE TABLE tactics (
    team_id TEXT NOT NULL,
    defense INT NOT NULL,
    center INT NOT NULL,
    attack INT NOT NULL,

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

    extra_attack_1  INT DEFAULT 0,
    extra_attack_2  INT DEFAULT 0,
    extra_attack_3  INT DEFAULT 0,
    extra_attack_4  INT DEFAULT 0,
    extra_attack_5  INT DEFAULT 0,
    extra_attack_6  INT DEFAULT 0,
    extra_attack_7  INT DEFAULT 0,
    extra_attack_8  INT DEFAULT 0,
    extra_attack_9  INT DEFAULT 0,
    extra_attack_10 INT DEFAULT 0,

    PRIMARY KEY(team_id)
)
