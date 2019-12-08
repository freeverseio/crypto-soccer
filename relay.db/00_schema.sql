SET TIMEZONE TO 'UTC';

-- CREATE TABLE verses (
--     verse BIGINT NOT NULL,
--     timestamp TIME
-- )

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
);

CREATE TABLE trainings (
    team_id TEXT NOT NULL,
    special_player_shirt INT CHECK (special_player_shirt >= -1 AND special_player_shirt <= 24),

    goalkeepers_defence INT CHECK (goalkeepers_defence >= 0),
    goalkeepers_speed INT CHECK (goalkeepers_speed >= 0),
    goalkeepers_pass INT CHECK (goalkeepers_pass >= 0),
    goalkeepers_shoot INT CHECK (goalkeepers_shoot >= 0),
    goalkeepers_endurance INT CHECK (goalkeepers_endurance >= 0),

    defenders_defence INT CHECK (defenders_defence >= 0),
    defenders_speed INT CHECK (defenders_speed >= 0),
    defenders_pass INT CHECK (defenders_pass >= 0),
    defenders_shoot INT CHECK (defenders_shoot >= 0),
    defenders_endurance INT CHECK (defenders_endurance >= 0),

    midfielders_defence INT CHECK (midfielders_defence >= 0),
    midfielders_speed INT CHECK (midfielders_speed >= 0),
    midfielders_pass INT CHECK (midfielders_pass >= 0),
    midfielders_shoot INT CHECK (midfielders_shoot >= 0),
    midfielders_endurance INT CHECK (midfielders_endurance >= 0),

    attackers_defence INT CHECK (attackers_defence >= 0),
    attackers_speed INT CHECK (attackers_speed >= 0),
    attackers_pass INT CHECK (attackers_pass >= 0),
    attackers_shoot INT CHECK (attackers_shoot >= 0),
    attackers_endurance INT CHECK (attackers_endurance >= 0),

    special_player_defence INT CHECK (special_player_defence >= 0),
    special_player_speed INT CHECK (special_player_speed >= 0),
    special_player_pass INT CHECK (special_player_pass >= 0),
    special_player_shoot INT CHECK (special_player_shoot >= 0),
    special_player_endurance INT CHECK (special_player_endurance >= 0),

    PRIMARY KEY(team_id)
);

CREATE TABLE create_special_players (
    uuid UUID NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    player_id TEXT NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY (uuid)
);

CREATE TABLE destroy_special_players (
    uuid UUID NOT NULL REFERENCES create_special_players(uuid),
    timestamp TIMESTAMP NOT NULL
);
