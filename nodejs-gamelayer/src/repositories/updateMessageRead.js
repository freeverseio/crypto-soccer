const PostgresSQLService = require('../services/PostgresSQLService');

const updateMessageReadQuery = {
  text: `
    UPDATE
      inbox
    SET
      is_read=true
    WHERE
      id=$1
    `,
};

const updateMessageRead = async ({ id }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [id];

  try {
    return await pool.query(updateMessageReadQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updateMessageRead;
