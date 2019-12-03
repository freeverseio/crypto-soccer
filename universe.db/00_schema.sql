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
    team_id NUMERIC(78,0) NOT NULL,
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
    PRIMARY KEY(team_id),
    FOREIGN KEY (timezone_idx, country_idx) REFERENCES countries(timezone_idx, country_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);

CREATE TABLE players (
    player_id NUMERIC(78,0) NOT NULL,
    name TEXT NOT NULL,
    team_id NUMERIC(78,0) NOT NULL REFERENCES teams(team_id),
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
    frozen BOOLEAN NOT NULL DEFAULT FALSE,
    red_card_matches_left INT NOT NULL DEFAULT 0,
    injury_matches_left INT NOT NULL DEFAULT 0,
    PRIMARY KEY(player_id)
);

CREATE TABLE matches (
    timezone_idx INT NOT NULL,
    country_idx INT NOT NULL,
    league_idx INT NOT NULL,
    match_day_idx INT NOT NULL,
    match_idx INT NOT NULL,
    home_team_id NUMERIC(78,0) REFERENCES teams(team_id),
    visitor_team_id NUMERIC(78,0) REFERENCES teams(team_id),
    home_goals INT,
    visitor_goals INT,
    home_match_log TEXT DEFAULT '0',
    visitor_match_log TEXT DEFAULT '0',
    PRIMARY KEY(timezone_idx,country_idx, league_idx, match_day_idx, match_idx),
    FOREIGN KEY (timezone_idx, country_idx, league_idx) REFERENCES leagues(timezone_idx, country_idx, league_idx)
);

CREATE TABLE auction_states(
    state TEXT NOT NULL PRIMARY KEY
);
INSERT INTO auction_states(state) VALUES ('STARTED');
INSERT INTO auction_states(state) VALUES ('ASSET_FROZEN');
INSERT INTO auction_states(state) VALUES ('PAYING');
INSERT INTO auction_states(state) VALUES ('PAID');
INSERT INTO auction_states(state) VALUES ('WITHDRAWAL');
INSERT INTO auction_states(state) VALUES ('NO_BIDS');
INSERT INTO auction_states(state) VALUES ('CANCELLED_BY_SELLER');
INSERT INTO auction_states(state) VALUES ('FAILED');

CREATE TABLE auctions (
    uuid UUID NOT NULL,
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price INT NOT NULL,
    rnd INT NOT NULL,
    valid_until TEXT NOT NULL,
    signature TEXT NOT NULL,
    state TEXT NOT NULL REFERENCES auction_states(state),
    state_extra TEXT NOT NULL DEFAULT '',
    payment_url TEXT NOT NULL DEFAULT '',
    seller TEXT NOT NULL,
    PRIMARY KEY(uuid)
);

CREATE TABLE bid_states(
    state TEXT NOT NULL PRIMARY KEY
);
INSERT INTO bid_states(state) VALUES ('ACCEPTED');
INSERT INTO bid_states(state) VALUES ('REFUSED');
INSERT INTO bid_states(state) VALUES ('PAYING');
INSERT INTO bid_states(state) VALUES ('PAID');
INSERT INTO bid_states(state) VALUES ('FAILED');

CREATE TABLE bids (
    auction UUID NOT NULL REFERENCES auctions(uuid),
    extra_price INT NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state TEXT NOT NULL REFERENCES bid_states(state),
    state_extra TEXT NOT NULL DEFAULT '',
    payment_id TEXT NOT NULL DEFAULT '',
    payment_url TEXT NOT NULL DEFAULT '',
    payment_deadline TEXT NOT NULL DEFAULT '0',
    PRIMARY KEY(auction, extra_price)
);

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

CREATE TABLE training (
    timestamp BIGINT NOT NULL,
    team_id TEXT NOT NULL,
    special_player_shirt INT CHECK (special_player_shirt >= -1 AND special_player_shirt <= 24),

    goalkeepers_defence INT CHECK (goalkeepers_defence >= 0),
    goalkeepers_speed INT CHECK (goalkeepers_speed >= 0),
    goalkeepers_pass INT CHECK (goalkeepers_pass >= 0),
    goalkeepers_shoot INT CHECK (goalkeepers_shoot >= 0),
    goalkeepers_endurance INT CHECK (goalkeepers_endurance >= 0),

    defencers_defencer INT CHECK (defencers_defencer >= 0),
    defencers_speed INT CHECK (defencers_speed >= 0),
    defencers_pass INT CHECK (defencers_pass >= 0),
    defencers_shoot INT CHECK (defencers_shoot >= 0),
    defencers_endurance INT CHECK (defencers_endurance >= 0),

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

    PRIMARY KEY(timestamp, team_id)
);


