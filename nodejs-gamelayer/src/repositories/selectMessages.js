const PostgresSQLService = require('../services/PostgresSQLService');

const selectMessagesQuery = ({ auctionId }) => {
  return {
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
      ${auctionId ? 'AND auction_id = $5' : ''}
      ORDER BY created_at DESC
    LIMIT $3
    OFFSET $4
  `,
  };
};

const selectMessages = async ({ destinatary, auctionId, createdAt, offset, limit }) => {
  const pool = await PostgresSQLService.getPool();
  const values = auctionId
    ? [destinatary, createdAt, limit, offset, auctionId]
    : [destinatary, createdAt, limit, offset];

  try {
    const { rows } = await pool.query(selectMessagesQuery({ auctionId }), values);
    return rows;
  } catch (e) {
    throw e;
  }
};

module.exports = selectMessages;
