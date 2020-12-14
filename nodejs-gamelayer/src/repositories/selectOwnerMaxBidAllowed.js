const PostgresSQLService = require('../services/PostgresSQLService');

const selectOwnerMaxBidAllowedQuery = {
  name: 'owner-max-bid-allowed-by-owner',
  text: `
    SELECT
      max_bid_allowed
    FROM
      owner_props
    WHERE
      owner = $1
  `,
};

const selectOwnerMaxBidAllowed = async ({ owner }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [owner];

  try {
    const { rows } = await pool.query(selectOwnerMaxBidAllowedQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectOwnerMaxBidAllowed;
