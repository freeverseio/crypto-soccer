const PostgresSQLService = require('../services/PostgresSQLService');

const selectPlayerNameQuery = {
  name: 'player-name-by-player-id',
  text: `
    SELECT
      player_name
    FROM
      player_props
    WHERE
      player_id = $1
  `,
};

const selectPlayerName = async ({ playerId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [playerId];

  try {
    const { rows } = await pool.query(selectPlayerNameQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectPlayerName;
