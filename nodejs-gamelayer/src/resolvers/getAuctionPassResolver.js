const { selectTeamAuctionPass } = require('../repositories');

const getAuctionPassResolver = async (parent, { teamId }, context, info, schema) => {
  const teamAuctionPass = await selectTeamAuctionPass({ teamId });

  return teamAuctionPass;
};

module.exports = getAuctionPassResolver;
