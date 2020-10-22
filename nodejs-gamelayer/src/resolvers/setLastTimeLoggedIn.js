const { upsertTeamLastTimeLoggedIn } = require('../repositories');

const setLastTimeLoggedIn = async (_, { teamId }) => {
  try {
    await upsertTeamLastTimeLoggedIn({ teamId });

    return true;
  } catch (e) {
    return e;
  }
};

module.exports = setLastTimeLoggedIn;
