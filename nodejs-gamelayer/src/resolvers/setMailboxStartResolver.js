const { insertTeamMailboxStartedAt } = require('../repositories');

const setMailboxStartResolver = async (_, { teamId }) => {
  try {
    await insertTeamMailboxStartedAt({ teamId });

    return true;
  } catch (e) {
    return e;
  }
};

module.exports = setMailboxStartResolver;
