const PostgresSQLService = require('../services/PostgresSQLService');

const selectMessagesQuery = {
  text: `
    SELECT
      id,
      destinatary,
      category,
      auction_id as auctionId,
      text_message as text,
      custom_image_url as customImageUrl,
      metadata,
      is_read as isRead
    FROM
      inbox
    WHERE
      destinatary = $1
      AND created_at >= $2
      AND id > $3
    LIMIT $4
  `,
};

const selectMessages = async ({ destinatary, createdAt, after, limit }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [destinatary, createdAt, after, limit];

  try {
    const { rows } = await pool.query(selectMessagesQuery, values);
    return rows;
  } catch (e) {
    throw e;
  }
};

module.exports = selectMessages;
