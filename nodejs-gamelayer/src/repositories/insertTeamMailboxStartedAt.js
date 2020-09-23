const PostgresSQLService = require('../services/PostgresSQLService');

const insertMailboxStartedAtQuery = {
  text: `
    INSERT INTO 
        team_props(
            team_id,
            mailbox_started_at
        )
    VALUES ($1, CURRENT_TIMESTAMP)
    `,
};

const insertMailboxStartedAt = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    return await pool.query(insertMailboxStartedAtQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = insertMailboxStartedAt;
