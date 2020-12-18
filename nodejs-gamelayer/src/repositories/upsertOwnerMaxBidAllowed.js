const PostgresSQLService = require('../services/PostgresSQLService');

const upsertMaxBidAlloweddQuery = {
  text: `
    INSERT INTO 
        owner_props(
            owner,
            max_bid_allowed
        )
    VALUES ($1, $2)
    ON CONFLICT (owner) DO UPDATE
    SET
    max_bid_allowed = $2
    `,
};

const upsertOwnerMaxBidAllowed = async ({ owner, maxBidAllowed }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [owner, maxBidAllowed];

  try {
    return await pool.query(upsertMaxBidAlloweddQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = upsertOwnerMaxBidAllowed;
