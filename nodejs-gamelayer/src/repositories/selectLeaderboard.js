const PostgresSQLService = require('../services/PostgresSQLService');

const selectLeaderboardQuery = {
  name: 'select-leadearboard-by-timezone-country-league',
  text: `
    SELECT
      leaderboard
    FROM
      league_props
    WHERE
      timezone_idx = $1
      AND country_idx = $2
      AND league_idx = $3
  `,
};

const selectLeaderboard = async ({ timezoneIdx, countryIdx, leagueIdx }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [timezoneIdx, countryIdx, leagueIdx];

  try {
    const { rows } = await pool.query(selectLeaderboardQuery, values);
    return rows[0];
  } catch (e) {
    throw e;
  }
};

module.exports = selectLeaderboard;
