const { selectPlayerName } = require('../repositories');

const secondaryPlayerByPlayerIdResolver = async (parent, args, context, info, schema) => {
  if (!parent.secondaryPlayerId) {
    return;
  }
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'playerByPlayerId',
    args: {
      playerId: parent.secondaryPlayerId,
    },
    context,
    info,
  });

  if (!result) {
    return;
  }

  const playerName = await selectPlayerName({ playerId: parent.secondaryPlayerId });

  result.name = playerName && playerName.player_name ? playerName.player_name : result.name;
  return result;
};

module.exports = secondaryPlayerByPlayerIdResolver;
