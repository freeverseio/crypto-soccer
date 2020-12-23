const PostgresSQLService = require('../services/PostgresSQLService');

const upsertCachedQueryDataQuery = {
  text: `
    INSERT INTO 
        query_cache(
            key,
            data,
            updated_at
        )
    VALUES ($1, $2, CURRENT_TIMESTAMP)
    ON CONFLICT (key) DO UPDATE
    SET
      data = $2,
      updated_at = CURRENT_TIMESTAMP
    `,
};

const upsertCachedQueryData = async ({ key, data }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [key, data];

  try {
    return await pool.query(upsertCachedQueryDataQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = upsertCachedQueryData;
