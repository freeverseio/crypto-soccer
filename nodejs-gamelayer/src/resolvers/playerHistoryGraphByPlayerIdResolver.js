const HorizonService = require('../services/HorizonService.js');

const playerHistoryGraphByPlayerIdResolver = async (parent, args, context, info, schema) => {
  const playerId = parent.playerId;
  const first = args.first;
  const step = 28;

  const playerHistory = await HorizonService.getPlayerHistory({playerId, count: (first * step)})
  console.log(playerHistory.nodes.length)

  const playerHistoryGraph = [];
  for (let i = 0 ; i < playerHistory.nodes.length ; i += step) {
    playerHistoryGraph.push(playerHistory.nodes[i]);
  }
  return { nodes: playerHistoryGraph };
};

module.exports = playerHistoryGraphByPlayerIdResolver;
