const PostgresSQLService = require('../services/PostgresSQLService');

const selectNumUnreadMessagesQuery = {
  text: `
    SELECT
      count(id) as num
    FROM
      inbox
    WHERE
      destinatary = $1
      AND is_read = false
      AND created_at >= $2
  `,
};

const selectNumUnreadMessages = async ({ destinatary, createdAt }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [destinatary, createdAt];

  try {
    const { rows } = await pool.query(selectNumUnreadMessagesQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectNumUnreadMessages;
