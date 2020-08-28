const PostgresSQLService = require('../services/PostgresSQLService');

const upsertLeaderboardQuery = {
  text: `
    INSERT INTO 
        league_props(
            timezone_idx,
            country_idx,
            league_idx,
            leaderboard
        )
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (timezone_idx, country_idx, league_idx) DO UPDATE
    SET
        leaderboard = $4
    `,
};

const upsertLeaderboard = async ({
  countryIdx,
  leagueIdx,
  timezoneIdx,
  leaderboard,
}) => {
  const pool = await PostgresSQLService.getPool();
  const values = [countryIdx, leagueIdx, timezoneIdx, leaderboard];

  try {
    return await pool.query(upsertLeaderboardQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = upsertLeaderboard;
