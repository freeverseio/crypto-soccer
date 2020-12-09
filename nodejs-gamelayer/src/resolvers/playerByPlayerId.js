const { selectPlayerName } = require('../repositories');

const playerByPlayerId = async (parent, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'playerByPlayerId',
    args: {
      playerId: parent.playerId,
    },
    context,
    info,
  });

  if (!result) {
    return;
  }

  const playerName = await selectPlayerName({ playerId: parent.playerId });

  result.name = playerName && playerName.player_name ? playerName.player_name : result.name;
  return result;
};

module.exports = playerByPlayerId;
