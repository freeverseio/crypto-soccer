const PostgresSQLService = require('../services/PostgresSQLService');

const selectTeamMaximumBidQuery = {
  name: 'team-allowed-bid-by-team-id',
  text: `
    SELECT
      maximum_bid
    FROM
      team_props
    WHERE
      team_id = $1
  `,
};

const selectTeamMaximumBid = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    const { rows } = await pool.query(selectTeamMaximumBidQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectTeamMaximumBid;
