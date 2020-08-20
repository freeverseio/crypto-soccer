const PostgresSQLService = require('../services/PostgresSQLService');

const updateTeamMaximumBidQuery = {
    text: `
    INSERT INTO 
        team_props(
            team_id,
            maximum_bid
        )
    VALUES ($1, $2)
    ON CONFLICT (team_id) DO UPDATE
    SET
        team_name = $2
    `,
  };

const updateTeamMaximumBid = async ({ teamId, teamMaximumBid }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId, teamMaximumBid];

  try {
    return await pool.query(updateTeamMaximumBidQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateTeamMaximumBid;
