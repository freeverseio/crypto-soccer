const PostgresSQLService = require('../services/PostgresSQLService');

const updateLastCheckedQuery = {
  text: `
    UPDATE
      mailbox_cron
    SET
      last_time_checked = $1
    WHERE
      entity_key = $2
    `,
};

const updateLastChecked = async ({ entity, lastChecked }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [lastChecked, entity];

  try {
    return await pool.query(updateLastCheckedQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateLastChecked;
