const PostgresSQLService = require('../services/PostgresSQLService');

const selectMessagesQuery = {
  text: `
    SELECT
      id,
      destinatary,
      category,
      auction_id,
      title,
      text_message as text,
      custom_image_url,
      metadata::TEXT,
      is_read,
      created_at
    FROM
      inbox
    WHERE
      destinatary = $1
      AND created_at >= $2
    LIMIT $3
    OFFSET $4
  `,
};

const selectMessages = async ({ destinatary, createdAt, offset, limit }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [destinatary, createdAt, limit, offset];

  try {
    const { rows } = await pool.query(selectMessagesQuery, values);
    return rows;
  } catch (e) {
    throw e;
  }
};

module.exports = selectMessages;
