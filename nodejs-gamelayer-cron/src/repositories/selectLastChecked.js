const PostgresSQLService = require('../services/PostgresSQLService');

const selectLastCheckedQuery = {
  text: `
    SELECT
      last_time_checked
    FROM
      mailbox_cron
    WHERE
      entity_key = $1
  `,
};

const selectLastChecked = async ({ entity }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [entity];

  try {
    const { rows } = await pool.query(selectLastCheckedQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectLastChecked;
