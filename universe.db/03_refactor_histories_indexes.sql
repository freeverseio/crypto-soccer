ALTER TABLE players_histories DROP CONSTRAINT players_histories_pkey;
ALTER TABLE players_histories ADD PRIMARY KEY(player_id, block_number);

ALTER TABLE teams_histories DROP CONSTRAINT teams_histories_pkey;
ALTER TABLE teams_histories ADD PRIMARY KEY(team_id, block_number);

ALTER TABLE tactics_histories DROP CONSTRAINT tactics_histories_pkey;
ALTER TABLE tactics_histories ADD PRIMARY KEY(team_id, block_number);

ALTER TABLE trainings_histories DROP CONSTRAINT trainings_histories_pkey;
ALTER TABLE trainings_histories ADD PRIMARY KEY(team_id, block_number);
