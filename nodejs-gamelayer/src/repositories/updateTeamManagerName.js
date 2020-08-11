const PostgresSQLService = require('../services/PostgresSQLService');

const updateTeamManagerNameQuery = {
    text: `
    INSERT INTO 
        team_props(
            team_id,
            team_manager_name
        )
    VALUES ($1, $2)
    ON CONFLICT (team_id) DO UPDATE
    SET
        team_manager_name = $2
    `,
  };

const updateTeamManagerName = async ({ teamId, teamManagerName }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId, teamManagerName];

  try {
    return await pool.query(updateTeamManagerNameQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateTeamManagerName;
