const HorizonService = require('../services/HorizonService.js');

const playerHistoryGraphByPlayerIdResolver = async (parent, args, context, info, schema) => {
  const teamId = parent.teamId;
  const playerId = parent.playerId;
  const first = args.first;

  const matchesBlockNumbers = await HorizonService.getMatchesPlayedByTeamId({ teamId });
  const playerHistoryGraph = [];
  for (const blockNumber of matchesBlockNumbers) {
    const encodedSkills = await HorizonService.getEncodedSkillsByBlockNumberPlayerId({
      playerId,
      blockNumber: blockNumber.blockNumber,
    });

    playerHistoryGraph.push({ encodedSkills: encodedSkills });
  }

  return { nodes: playerHistoryGraph.slice(0, first) };
};

module.exports = playerHistoryGraphByPlayerIdResolver;
