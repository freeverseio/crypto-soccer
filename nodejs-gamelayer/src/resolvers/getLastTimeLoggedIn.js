const dayjs = require('dayjs');
const { selectTeamLastTimeLoggedIn } = require('../repositories');

const getLastTimeLoggedIn = async (parent, { teamId }, context, info, schema) => {
  const lastTimeLoggedIn = await selectTeamLastTimeLoggedIn({ teamId });

  return dayjs(lastTimeLoggedIn).format();
};

module.exports = getLastTimeLoggedIn;
