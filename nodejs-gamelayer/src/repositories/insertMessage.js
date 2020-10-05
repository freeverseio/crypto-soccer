const PostgresSQLService = require('../services/PostgresSQLService');

const insertMessageQuery = {
  text: `
    INSERT INTO 
        inbox(
          destinatary,
          category,
          auction_id,
          text_message,
          custom_image_url,
          metadata,
          is_read
        )
    VALUES ($1, $2, $3, $4, $5, $6, false)
    RETURNING id
    `,
};

const insertMessage = async ({ destinatary, category, auctionId, text, customImageUrl, metadata }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [destinatary, category, auctionId, text, customImageUrl, metadata];

  try {
    const { rows } = await pool.query(insertMessageQuery, values);
    const { id } = rows[0];
    return id;
  } catch (e) {
    throw e;
  }
};

module.exports = insertMessage;
