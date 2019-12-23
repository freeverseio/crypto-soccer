CREATE TABLE tactics (
    verse BIGINT NOT NULL DEFAULT 0,
    team_id TEXT NOT NULL,

    tactic_id INT NOT NULL,

    shirt_0 INT NOT NULL CHECK (shirt_0 >= 0),
    shirt_1 INT NOT NULL CHECK (shirt_1 >= 0),
    shirt_2 INT NOT NULL CHECK (shirt_2 >= 0),
    shirt_3 INT NOT NULL CHECK (shirt_3 >= 0),
    shirt_4 INT NOT NULL CHECK (shirt_4 >= 0),
    shirt_5 INT NOT NULL CHECK (shirt_5 >= 0),
    shirt_6 INT NOT NULL CHECK (shirt_6 >= 0),
    shirt_7 INT NOT NULL CHECK (shirt_7 >= 0),
    shirt_8 INT NOT NULL CHECK (shirt_8 >= 0),
    shirt_9 INT NOT NULL CHECK (shirt_9 >= 0),
    shirt_10 INT NOT NULL CHECK (shirt_10 >= 0),

    shirt_11 INT NOT NULL CHECK (shirt_11 >= 0),
    shirt_12 INT NOT NULL CHECK (shirt_12 >= 0),
    shirt_13 INT NOT NULL CHECK (shirt_13 >= 0),

    extra_attack_1  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_2  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_3  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_4  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_5  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_6  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_7  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_8  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_9  BOOLEAN NOT NULL DEFAULT FALSE,
    extra_attack_10 BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (verse, team_id)
);

CREATE VIEW current_tactic AS SELECT * FROM tactics WHERE verse=0;

CREATE TABLE trainings (
    verse BIGINT NOT NULL DEFAULT 0,
    team_id TEXT NOT NULL,

    special_player_shirt INT NOT NULL CHECK (special_player_shirt >= -1 AND special_player_shirt <= 24),

    goalkeepers_defence INT NOT NULL DEFAULT 0 CHECK (goalkeepers_defence >= 0),
    goalkeepers_speed INT NOT NULL DEFAULT 0 CHECK (goalkeepers_speed >= 0),
    goalkeepers_pass INT NOT NULL DEFAULT 0 CHECK (goalkeepers_pass >= 0),
    goalkeepers_shoot INT NOT NULL DEFAULT 0 CHECK (goalkeepers_shoot >= 0),
    goalkeepers_endurance INT NOT NULL DEFAULT 0 CHECK (goalkeepers_endurance >= 0),

    defenders_defence INT NOT NULL DEFAULT 0 CHECK (defenders_defence >= 0),
    defenders_speed INT NOT NULL DEFAULT 0 CHECK (defenders_speed >= 0),
    defenders_pass INT NOT NULL DEFAULT 0 CHECK (defenders_pass >= 0),
    defenders_shoot INT NOT NULL DEFAULT 0 CHECK (defenders_shoot >= 0),
    defenders_endurance INT NOT NULL DEFAULT 0 CHECK (defenders_endurance >= 0),

    midfielders_defence INT NOT NULL DEFAULT 0 CHECK (midfielders_defence >= 0),
    midfielders_speed INT NOT NULL DEFAULT 0 CHECK (midfielders_speed >= 0),
    midfielders_pass INT NOT NULL DEFAULT 0 CHECK (midfielders_pass >= 0),
    midfielders_shoot INT NOT NULL DEFAULT 0 CHECK (midfielders_shoot >= 0),
    midfielders_endurance INT NOT NULL DEFAULT 0 CHECK (midfielders_endurance >= 0),

    attackers_defence INT NOT NULL DEFAULT 0 CHECK (attackers_defence >= 0),
    attackers_speed INT NOT NULL DEFAULT 0 CHECK (attackers_speed >= 0),
    attackers_pass INT NOT NULL DEFAULT 0 CHECK (attackers_pass >= 0),
    attackers_shoot INT NOT NULL DEFAULT 0 CHECK (attackers_shoot >= 0),
    attackers_endurance INT NOT NULL DEFAULT 0 CHECK (attackers_endurance >= 0),

    special_player_defence INT NOT NULL DEFAULT 0 CHECK (special_player_defence >= 0),
    special_player_speed INT NOT NULL DEFAULT 0 CHECK (special_player_speed >= 0),
    special_player_pass INT NOT NULL DEFAULT 0 CHECK (special_player_pass >= 0),
    special_player_shoot INT NOT NULL DEFAULT 0 CHECK (special_player_shoot >= 0),
    special_player_endurance INT NOT NULL DEFAULT 0 CHECK (special_player_endurance >= 0),

    PRIMARY KEY (verse, team_id)
 );


CREATE VIEW current_trainings AS SELECT * FROM trainings WHERE verse=0;