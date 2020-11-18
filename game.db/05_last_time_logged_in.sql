ALTER TABLE team_props ALTER COLUMN team_name DROP NOT NULL;
ALTER TABLE team_props ALTER COLUMN team_manager_name DROP NOT NULL;
ALTER TABLE team_props ADD COLUMN last_time_logged_in timestamp without time zone;