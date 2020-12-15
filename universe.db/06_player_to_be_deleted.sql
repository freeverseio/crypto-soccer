ALTER TABLE players ADD COLUMN voided BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX index_players_voided ON players(voided);