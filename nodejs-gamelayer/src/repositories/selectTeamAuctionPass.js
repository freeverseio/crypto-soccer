const PostgresSQLService = require('../services/PostgresSQLService');

const selectTeamAuctionPassQuery = {
  name: 'team-last-time-logged-in',
  text: `
    SELECT
      team_id,
      auction_pass
    FROM
      team_props
    WHERE
      team_id = $1
  `,
};

const selectTeamAuctionPass = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    const { rows } = await pool.query(selectTeamAuctionPassQuery, values);
    if (rows[0]) {
      return rows[0];
    }

    return '';
  } catch (e) {
    throw e;
  }
};

module.exports = selectTeamAuctionPass;
