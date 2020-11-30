const PostgresSQLService = require('../services/PostgresSQLService');

const upsertGetSocialIdQuery = {
  text: `
    INSERT INTO 
        team_props(
            team_id,
            get_social_id
        )
    VALUES ($1, $2)
    ON CONFLICT (team_id) DO UPDATE
    SET
      get_social_id = $2
    `,
};

const upsertTeamGetSocialId = async ({ teamId, getSocialId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId, getSocialId];

  try {
    return await pool.query(upsertGetSocialIdQuery, values);
  } catch (e) {
    throw e;
  }
};

module.exports = upsertTeamGetSocialId;
