const PostgresSQLService = require('../services/PostgresSQLService');

const updatePlayerNameQuery = {
    text: `
    INSERT INTO player_props (
        player_id,
        player_name
    )
    VALUES ($1, $2)
    ON CONFLICT (player_id) DO UPDATE
    SET
        player_name = $2
    `,
  };

const updatePlayerName = async ({ playerId, playerName }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [playerId, playerName];

  try {
    return await pool.query(updatePlayerNameQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = updatePlayerName;
