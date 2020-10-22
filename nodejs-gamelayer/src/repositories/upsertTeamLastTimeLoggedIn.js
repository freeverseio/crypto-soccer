const PostgresSQLService = require('../services/PostgresSQLService');

const upsertLastTimeLoggedInQuery = {
  text: `
    INSERT INTO 
        team_props(
            team_id,
            last_time_logged_in
        )
    VALUES ($1, CURRENT_TIMESTAMP)
    ON CONFLICT (team_id) DO UPDATE
    SET
      last_time_logged_in = CURRENT_TIMESTAMP
    `,
};

const upsertLastTimeLoggedIn = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    return await pool.query(upsertLastTimeLoggedInQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = upsertLastTimeLoggedIn;
