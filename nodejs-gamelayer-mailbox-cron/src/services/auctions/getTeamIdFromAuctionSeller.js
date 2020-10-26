const HorizonService = require('../HorizonService');

const getTeamIdFromAuctionSeller = async ({ auction }) => {
  const { seller, playerId } = auction;
  const possibleTeamIds = await HorizonService.getTeamIdsFromOwner({
    owner: seller,
  });
  const cleanPossibleTeamIds = possibleTeamIds.map((r) => r.teamId);

  if (possibleTeamIds.length == 1) {
    return cleanPossibleTeamIds[0];
  }

  const recentTeamsForPlayer = await HorizonService.getPlayerHistoriesLast30BlockNumberTeams(
    { playerId }
  );
  const cleanRecentTeamsForPlayer = recentTeamsForPlayer.map((r) => r.teamId);

  for (const teamId of cleanRecentTeamsForPlayer) {
    const playerSellerTeamId = cleanPossibleTeamIds.find((t) => t == teamId);
    if (playerSellerTeamId) {
      return playerSellerTeamId;
    }
  }
};

module.exports = getTeamIdFromAuctionSeller;
