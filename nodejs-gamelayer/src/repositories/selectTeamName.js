const PostgresSQLService = require('../services/PostgresSQLService');

const selectTeamNameQuery = {
  name: 'team-name-by-team-id',
  text: `
    SELECT
      team_name
    FROM
      team_props
    WHERE
      team_id = $1
  `,
};

const selectTeamName = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    const { rows } = await pool.query(selectTeamNameQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectTeamName;
