const PostgresSQLService = require('../services/PostgresSQLService');

const updateTeamAssignedAtQuery = {
  text: `
    INSERT INTO 
        team_props(
            team_id,
            assigned_at
        )
    VALUES ($1, $2)
    ON CONFLICT (team_id) DO UPDATE
    SET
    assigned_at = $2
    `,
};

const updateTeamAssignedAt = async ({ teamId, assignedAt }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId, assignedAt];

  try {
    return await pool.query(updateTeamAssignedAtQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateTeamAssignedAt;
