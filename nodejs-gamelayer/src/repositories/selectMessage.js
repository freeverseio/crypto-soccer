const PostgresSQLService = require('../services/PostgresSQLService');

const selectMessagesQuery = {
  text: `
    SELECT
      destinatary,
      category,
      auction_id as auctionId,
      text_message as text,
      custom_image_url as customImageUrl,
      metadata
    FROM
      inbox
    WHERE
      destinatary = $1
      AND created_at >= $2
  `,
};

const selectMessages = async ({ destinatary, createdAt }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [destinatary, createdAt];
  console.log('selectMessages -> values', values);

  try {
    const { rows } = await pool.query(selectMessagesQuery, values);
    return rows;
  } catch (e) {
    throw e;
  }
};

module.exports = selectMessages;
