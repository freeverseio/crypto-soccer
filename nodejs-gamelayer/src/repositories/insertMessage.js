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
          metadata
        )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id
    `,
};

const insertMessage = async ({
  destinatary,
  category,
  auction_id,
  text_message,
  custom_image_url,
  metadata,
}) => {
  const pool = await PostgresSQLService.getPool();
  const values = [
    destinatary,
    category,
    auction_id,
    text_message,
    custom_image_url,
    metadata,
  ];

  try {
    const { rows } = await pool.query(insertMessageQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = insertMessage;
