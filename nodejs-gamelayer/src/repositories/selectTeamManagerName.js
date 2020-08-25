const PostgresSQLService = require('../services/PostgresSQLService');

const selectTeamManagerNameQuery = {
  name: 'team-manager-name-by-team-id',
  text: `
    SELECT
      team_manager_name
    FROM
      team_props
    WHERE
      team_id = $1
  `,
};

const selectTeamManagerName = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    const { rows } = await pool.query(selectTeamManagerNameQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectTeamManagerName;
