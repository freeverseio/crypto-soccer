CREATE INDEX CONCURRENTLY index_teams_timezone_idx ON teams (timezone_idx);
CREATE INDEX CONCURRENTLY index_teams_country_idx ON teams (country_idx);
CREATE INDEX CONCURRENTLY index_teams_league_idx ON teams (league_idx);

CREATE INDEX CONCURRENTLY index_matches_timezone_idx ON matches (timezone_idx);
CREATE INDEX CONCURRENTLY index_matches_country_idx ON matches (country_idx);
CREATE INDEX CONCURRENTLY index_matches_league_idx ON matches (league_idx);
CREATE INDEX CONCURRENTLY index_matches_match_day_idx ON matches (match_day_idx);
CREATE INDEX CONCURRENTLY index_matches_match_idx ON matches (match_idx);

CREATE INDEX CONCURRENTLY index_leagues_league_idx ON leagues (league_idx);


CREATE INDEX CONCURRENTLY index_players_team_id ON players (team_id);
