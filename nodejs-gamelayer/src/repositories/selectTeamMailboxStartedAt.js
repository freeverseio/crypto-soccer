const PostgresSQLService = require('../services/PostgresSQLService');

const selectTeamMailboxStartedAtQuery = {
  name: 'team-mailbox-started-at-by-team-id',
  text: `
    SELECT
      mailbox_started_at
    FROM
      team_props
    WHERE
      team_id = $1
  `,
};

const selectTeamMailboxStartedAt = async ({ teamId }) => {
  const pool = await PostgresSQLService.getPool();
  const values = [teamId];

  try {
    const { rows } = await pool.query(selectTeamMailboxStartedAtQuery, values);
    if (rows[0] && rows[0].mailbox_started_at) {
      const { mailbox_started_at } = rows[0];
      return mailbox_started_at;
    }

    return '';
  } catch (e) {
    throw e;
  }
};

module.exports = selectTeamMailboxStartedAt;
