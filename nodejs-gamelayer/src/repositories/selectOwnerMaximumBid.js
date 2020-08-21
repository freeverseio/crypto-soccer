const PostgresSQLService = require('../services/PostgresSQLService');

const selectOwnerMaximumBidQuery = {
  name: 'maximum-allowed-bid-by-owner',
  text: `
    SELECT
      maximum_bid
    FROM
      owner_props
    WHERE
      owner = $1
  `,
};

const selectOwnerMaximumBid = async ({ owner }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [owner];

  try {
    const { rows } = await pool.query(selectOwnerMaximumBidQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectOwnerMaximumBid;
