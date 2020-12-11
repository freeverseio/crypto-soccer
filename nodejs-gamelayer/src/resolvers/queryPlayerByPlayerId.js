const { selectPlayerName } = require('../repositories');

const queryPlayerByPlayerId = async (parent, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'playerByPlayerId',
    args: {
      playerId: args.playerId,
    },
    context,
    info,
  });

  if (!result) {
    return;
  }

  const playerName = await selectPlayerName({ playerId: args.playerId });

  result.name = playerName && playerName.player_name ? playerName.player_name : result.name;
  return result;
};

module.exports = queryPlayerByPlayerId;
