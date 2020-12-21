const { selectCachedQueryByKey, upsertCachedQueryData } = require('../repositories');

const queryBestPlayers = async (parent, args, context, info, schema) => {
  const { limit } = args;
  const key = `get_best_players_limit_${limit}`;
  const cachedData = await selectCachedQueryByKey({ key });

  if (cachedData) {
    return cachedData;
  }

  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'getBestPlayers',
    args: {
      limit,
    },
    context,
    info,
  });

  if (!result) {
    return;
  }

  await upsertCachedQueryData({ key, data: result });

  return result;
};

module.exports = queryBestPlayers;
