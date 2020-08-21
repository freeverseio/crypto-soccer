const PostgresSQLService = require('../services/PostgresSQLService');

const updateOwnerMaximumBidQuery = {
  text: `
    INSERT INTO 
        owner_props(
            owner,
            maximum_bid
        )
    VALUES ($1, $2)
    ON CONFLICT (team_id) DO UPDATE
    SET
        maximum_bid = $2
    `,
};

const updateOwnerMaximumBid = async ({ owner, maximumBid }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [owner, maximumBid];

  try {
    return await pool.query(updateOwnerMaximumBidQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateOwnerMaximumBid;
