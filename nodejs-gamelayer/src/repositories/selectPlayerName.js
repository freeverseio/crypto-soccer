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
  console.log("ei")
  const pool = await PostgresSQLService.getPool();
  console.log("ei2")
  const values = [playerId];

  try {
    const { rows } = await pool.query(selectPlayerNameQuery, values);
    console.log("Da rows", rows)
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectPlayerName;
