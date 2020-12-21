const PostgresSQLService = require('../services/PostgresSQLService');

const selectCachedQueryByKeyQuery = {
  text: `
    SELECT
      data,
      updated_at
    FROM
      query_cache
    WHERE
      key = $1
  `,
};

const selectCachedQueryByKey = async ({ key }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [key];

  try {
    const { rows } = await pool.query(selectCachedQueryByKeyQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectCachedQueryByKey;
