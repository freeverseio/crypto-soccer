const PostgresSQLService = require('../services/PostgresSQLService');

const selectTeamLastTimeLoggedInQuery = {
  name: 'team-last-time-logged-in',
  text: `
    SELECT
      last_time_logged_in
    FROM
      team_props
    WHERE
      team_id = $1
  `,
};

const selectTeamLastTimeLoggedIn = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    const { rows } = await pool.query(selectTeamLastTimeLoggedInQuery, values);
    if (rows[0] && rows[0].last_time_logged_in) {
      const { last_time_logged_in } = rows[0];
      return last_time_logged_in;
    }

    return '';
  } catch (e) {
    throw e;
  }
};

module.exports = selectTeamLastTimeLoggedIn;
