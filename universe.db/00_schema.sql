CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);
comment on table params is E'@omit create,update,delete';

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE timezones (
    timezone_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx)
);
comment on table timezones is E'@omit create,update,delete';

CREATE TABLE countries (
    timezone_idx INT NOT NULL REFERENCES timezones(timezone_idx),
    country_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx, country_idx)
);
comment on table countries is E'@omit create,update,delete';

CREATE TABLE leagues (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    PRIMARY KEY(timezone_idx,country_idx, league_idx),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx)
);
comment on table leagues is E'@omit create,update,delete';

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
    match_log TEXT NOT NULL,
    PRIMARY KEY(team_id),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);
comment on table teams is E'@omit create,update,delete';

CREATE TABLE teams_histories (
    block_number BIGINT NOT NULL,
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
    match_log TEXT NOT NULL,
    PRIMARY KEY(block_number, team_id),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);
comment on table teams_histories is E'@omit create,update,delete';

CREATE TABLE players (
    name TEXT NOT NULL,
    player_id TEXT NOT NULL,
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
    tiredness INT NOT NULL,
    PRIMARY KEY(player_id)
);
comment on table players is E'@omit create,update,delete';

CREATE TABLE players_histories (
    player_id TEXT NOT NULL REFERENCES players(player_id),
    block_number BIGINT NOT NULL,
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
    tiredness INT NOT NULL,
    PRIMARY KEY(block_number, player_id)
);
comment on table players_histories is E'@omit create,update,delete';

CREATE TYPE match_state AS ENUM ('begin', 'half', 'end', 'cancelled');
CREATE TABLE matches (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    match_day_idx INT NOT NULL,
    match_idx INT NOT NULL,
    home_team_id TEXT REFERENCES teams(team_id),
    visitor_team_id TEXT REFERENCES teams(team_id),
    seed TEXT NOT NULL DEFAULT '',
    home_goals INT NOT NULL DEFAULT 0,
    visitor_goals INT NOT NULL DEFAULT 0,
    home_teamsumskills INT NOT NULL DEFAULT 0,
    visitor_teamsumskills INT NOT NULL DEFAULT 0,
    state match_state NOT NULL,
    state_extra TEXT NOT NULL DEFAULT '',
    start_epoch BIGINT NOT NULL,
    PRIMARY KEY(timezone_idx,country_idx, league_idx, match_day_idx, match_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);
comment on table matches is E'@omit create,update,delete';

CREATE TABLE matches_histories (
    block_number INT NOT NULL,
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    match_day_idx INT NOT NULL,
    match_idx INT NOT NULL,
    home_team_id TEXT REFERENCES teams(team_id),
    visitor_team_id TEXT REFERENCES teams(team_id),
    seed TEXT NOT NULL,
    home_goals INT NOT NULL,
    visitor_goals INT NOT NULL,
    home_teamsumskills INT NOT NULL,
    visitor_teamsumskills INT NOT NULL,
    state match_state NOT NULL,
    state_extra TEXT NOT NULL,
    start_epoch BIGINT NOT NULL,
    PRIMARY KEY(block_number, timezone_idx,country_idx, league_idx, match_day_idx, match_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);
comment on table matches_histories is E'@omit create,update,delete';

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
    primary_player_id TEXT REFERENCES players(player_id),
    secondary_player_id TEXT REFERENCES players(player_id),
    FOREIGN KEY (timezone_idx, country_idx, league_idx, match_day_idx, match_idx) REFERENCES matches(timezone_idx, country_idx, league_idx, match_day_idx, match_idx)
);
comment on table match_events is E'@omit create,update,delete';

CREATE TABLE verses (
    verse_number BIGINT NOT NULL,
    root TEXT NOT NULL,
    PRIMARY KEY (verse_number)
);
comment on table verses is E'@omit create,update,delete';



