const PostgresSQLService = require('../services/PostgresSQLService');

const updateTeamNameQuery = {
    text: `
    INSERT INTO 
        team_props(
            team_id,
            team_name,
            team_manager_name
        )
    VALUES ($1, $2, '')
    ON CONFLICT (team_id) DO UPDATE
    SET
        team_name = $2
    `,
  };

const updateTeamName = async ({ teamId, teamName }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId, teamName];

  try {
    return await pool.query(updateTeamNameQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateTeamName;
