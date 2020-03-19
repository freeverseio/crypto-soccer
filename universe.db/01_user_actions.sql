CREATE TABLE tactics (
    team_id TEXT NOT NULL REFERENCES teams(team_id),

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

    substitution_0_shirt INT NOT NULL CHECK (substitution_0_shirt >= 0 AND substitution_0_shirt <= 25),
    substitution_0_target INT NOT NULL CHECK (substitution_0_target >= 0 AND substitution_0_target <= 11),
    substitution_0_minute INT NOT NULL CHECK (substitution_0_minute >= 0 AND substitution_0_minute <= 90),

    substitution_1_shirt INT NOT NULL CHECK (substitution_1_shirt >= 0 AND substitution_1_shirt <= 25),
    substitution_1_target INT NOT NULL CHECK (substitution_1_target >= 0 AND substitution_1_target <= 11),
    substitution_1_minute INT NOT NULL CHECK (substitution_1_minute >= 0 AND substitution_1_minute <= 90),

    substitution_2_shirt INT NOT NULL CHECK (substitution_2_shirt >= 0 AND substitution_2_shirt <= 25),
    substitution_2_target INT NOT NULL CHECK (substitution_2_target >= 0 AND substitution_2_target <= 11),
    substitution_2_minute INT NOT NULL CHECK (substitution_2_minute >= 0 AND substitution_2_minute <= 90),

    extra_attack_1  BOOLEAN NOT NULL,
    extra_attack_2  BOOLEAN NOT NULL,
    extra_attack_3  BOOLEAN NOT NULL,
    extra_attack_4  BOOLEAN NOT NULL,
    extra_attack_5  BOOLEAN NOT NULL,
    extra_attack_6  BOOLEAN NOT NULL,
    extra_attack_7  BOOLEAN NOT NULL,
    extra_attack_8  BOOLEAN NOT NULL,
    extra_attack_9  BOOLEAN NOT NULL,
    extra_attack_10 BOOLEAN NOT NULL,

    PRIMARY KEY (team_id)
);
comment on table tactics is E'@omit create,delete';
comment on column tactics.team_id is E'@omit update';

CREATE TABLE tactics_histories (
    block_number BIGINT NOT NULL,
    team_id TEXT NOT NULL REFERENCES teams(team_id),

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

    substitution_0_shirt INT NOT NULL CHECK (substitution_0_shirt >= 0 AND substitution_0_shirt <= 25),
    substitution_0_target INT NOT NULL CHECK (substitution_0_target >= 0 AND substitution_0_target <= 11),
    substitution_0_minute INT NOT NULL CHECK (substitution_0_minute >= 0 AND substitution_0_minute <= 90),

    substitution_1_shirt INT NOT NULL CHECK (substitution_1_shirt >= 0 AND substitution_1_shirt <= 25),
    substitution_1_target INT NOT NULL CHECK (substitution_1_target >= 0 AND substitution_1_target <= 11),
    substitution_1_minute INT NOT NULL CHECK (substitution_1_minute >= 0 AND substitution_1_minute <= 90),

    substitution_2_shirt INT NOT NULL CHECK (substitution_2_shirt >= 0 AND substitution_2_shirt <= 25),
    substitution_2_target INT NOT NULL CHECK (substitution_2_target >= 0 AND substitution_2_target <= 11),
    substitution_2_minute INT NOT NULL CHECK (substitution_2_minute >= 0 AND substitution_2_minute <= 90),

    extra_attack_1  BOOLEAN NOT NULL,
    extra_attack_2  BOOLEAN NOT NULL,
    extra_attack_3  BOOLEAN NOT NULL,
    extra_attack_4  BOOLEAN NOT NULL,
    extra_attack_5  BOOLEAN NOT NULL,
    extra_attack_6  BOOLEAN NOT NULL,
    extra_attack_7  BOOLEAN NOT NULL,
    extra_attack_8  BOOLEAN NOT NULL,
    extra_attack_9  BOOLEAN NOT NULL,
    extra_attack_10 BOOLEAN NOT NULL,

    PRIMARY KEY (block_number, team_id)
);
comment on table tactics_histories is E'@omit create,update,delete';

CREATE TABLE trainings (
    team_id TEXT NOT NULL REFERENCES teams(team_id),

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

    PRIMARY KEY (team_id)
 );
comment on table trainings is E'@omit create,delete';
comment on column trainings.team_id is E'@omit update';

CREATE TABLE trainings_histories (
    block_number BIGINT NOT NULL,
    team_id TEXT NOT NULL REFERENCES teams(team_id),

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

    PRIMARY KEY (block_number, team_id)
);
comment on table trainings_histories is E'@omit create,update,delete';

CREATE TABLE teams_props (
    team_id TEXT NOT NULL REFERENCES teams(team_id),
    name TEXT NOT NULL,
    PRIMARY KEY (team_id)
);
comment on table teams_props is E'@omit create,delete';
comment on column teams_props.team_id is E'@omit update';

CREATE TABLE teams_props_histories (
    block_number BIGINT NOT NULL,
    team_id TEXT NOT NULL REFERENCES teams(team_id),
    name TEXT NOT NULL,
    PRIMARY KEY (block_number, team_id)
);
comment on table teams_props_histories is E'@omit create,update,delete'