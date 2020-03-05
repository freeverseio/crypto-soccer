CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE timezones (
    timezone_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx)
);

CREATE TABLE countries (
    timezone_idx INT NOT NULL REFERENCES timezones(timezone_idx),
    country_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx, country_idx)
);

CREATE TABLE leagues (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx,country_idx, league_idx),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx)
);

CREATE TABLE teams (
    team_id TEXT NOT NULL,
    name TEXT NOT NULL,
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    owner TEXT NOT NULL,
    league_idx INT NOT NULL,
    team_idx_in_league INT NOT NULL,
    points INT NOT NULL DEFAULT 0,
    w INT NOT NULL DEFAULT 0,
    d INT NOT NULL DEFAULT 0,
    l INT NOT NULL DEFAULT 0,
    goals_forward INT NOT NULL DEFAULT 0,
    goals_against INT NOT NULL DEFAULT 0,
    prev_perf_points TEXT NOT NULL DEFAULT '0',
    ranking_points TEXT NOT NULL DEFAULT '0',
    training_points INT NOT NULL DEFAULT 0,
    tactic TEXT NOT NULL DEFAULT '',
    PRIMARY KEY(team_id),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);

CREATE TABLE players (
    player_id TEXT NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY(player_id)
);

CREATE TABLE players_states (
    block_number BIGINT NOT NULL,
    player_id TEXT NOT NULL REFERENCES players(player_id),
    team_id TEXT NOT NULL REFERENCES teams(team_id),
    defence INT NOT NULL,
    speed INT NOT NULL,
    pass INT NOT NULL,
    shoot INT NOT NULL,
    endurance INT NOT NULL,
    shirt_number INT NOT NULL,
    preferred_position TEXT NOT NULL,
    potential INT NOT NULL, 
    day_of_birth INT NOT NULL, 
    encoded_skills TEXT NOT NULL,
    encoded_state TEXT NOT NULL,
    red_card BOOL NOT NULL DEFAULT FALSE,
    injury_matches_left INT NOT NULL DEFAULT 0,
    PRIMARY KEY(block_number, player_id)
);

CREATE VIEW current_players AS SELECT * FROM 
    players 
    LEFT JOIN LATERAL
    (SELECT * FROM players_states WHERE player_id = players.player_id ORDER BY block_number DESC LIMIT 1) o2 ON true;

CREATE TYPE match_state AS ENUM ('begin', 'half', 'end', 'cancel');
CREATE TABLE matches (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    match_day_idx INT NOT NULL,
    match_idx INT NOT NULL,
    home_team_id TEXT REFERENCES teams(team_id),
    visitor_team_id TEXT REFERENCES teams(team_id),
    home_goals INT NOT NULL DEFAULT 0,
    visitor_goals INT NOT NULL DEFAULT 0,
    home_match_log TEXT NOT NULL,
    visitor_match_log TEXT NOT NULL,
    state match_state NOT NULL,
    PRIMARY KEY(timezone_idx,country_idx, league_idx, match_day_idx, match_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);

CREATE TYPE match_event_type AS ENUM ('attack', 'yellow_card', 'red_card', 'injury_soft', 'injury_hard', 'substitution');
CREATE TABLE match_events (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    match_day_idx INT NOT NULL,
    match_idx INT NOT NULL,
    minute INT NOT NULL,
    type match_event_type NOT NULL,
    team_id TEXT NOT NULL REFERENCES teams(team_id),
    manage_to_shoot BOOLEAN NOT NULL DEFAULT 'false',
    is_goal BOOLEAN NOT NULL DEFAULT 'false',
    primary_player_id TEXT NOT NULL REFERENCES players(player_id),
    secondary_player_id TEXT REFERENCES players(player_id),
    FOREIGN KEY (timezone_idx, country_idx, league_idx, match_day_idx, match_idx) REFERENCES matches(timezone_idx, country_idx, league_idx, match_day_idx, match_idx)
);

-- CREATE TABLE teams_history (
--     teamId BIGINT NOT NULL REFERENCES teams(id),
--     blockNumber BIGINT NOT NULL,
--     inBlockIndex INT NOT NULL,
--     owner TEXT NOT NULL,
--     currentLeagueId BIGINT NOT NULL,
--     posInCurrentLeagueId INT NOT NULL,
--     prevLeagueId BIGINT NOT NULL,
--     posInPrevLeagueId INT NOT NULL,
--     PRIMARY KEY(teamId, blockNumber, inBlockIndex)
-- );

-- CREATE TABLE players (
--     id BIGINT NOT NULL,
--     monthOfBirthInUnixTime TEXT NOT NULL,
--     blockNumber BIGINT NOT NULL,
--     inBlockIndex INT NOT NULL,
--     teamId BIGINT NOT NULL REFERENCES teams(id),
--     state TEXT NOT NULL,
--     defence INT NOT NULL,
--     speed INT NOT NULL,
--     pass INT NOT NULL,
--     shoot INT NOT NULL,
--     endurance INT NOT NULL,
--     PRIMARY KEY(id)
-- );

-- CREATE TABLE players_history (
--     playerId BIGINT NOT NULL REFERENCES players(id),
--     blockNumber BIGINT NOT NULL,
--     inBlockIndex INT NOT NULL,
--     teamId BIGINT NOT NULL REFERENCES teams(id),
--     state TEXT NOT NULL,
--     defence INT NOT NULL,
--     speed INT NOT NULL,
--     pass INT NOT NULL,
--     shoot INT NOT NULL,
--     endurance INT NOT NULL,
--     PRIMARY KEY(playerId, blockNumber, inBlockIndex)
-- );

-- ********************************** USER ACTIONS ****************************************

CREATE TABLE tactics (
    verse BIGINT NOT NULL DEFAULT 0,
    timezone INT NOT NULL,
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

    PRIMARY KEY (verse, team_id)
);

CREATE VIEW upcoming_tactics AS SELECT * FROM tactics WHERE verse=9223372036854775807;

CREATE TABLE trainings (
    verse BIGINT NOT NULL DEFAULT 0,
    timezone INT NOT NULL,
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

    PRIMARY KEY (verse, team_id)
 );

 CREATE VIEW upcoming_trainings AS SELECT * FROM trainings WHERE verse=9223372036854775807;
