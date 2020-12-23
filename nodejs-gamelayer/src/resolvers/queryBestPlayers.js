const utc = require('dayjs/plugin/utc');
const dayjs = require('dayjs');
dayjs.extend(utc);
const { selectCachedQueryByKey, upsertCachedQueryData } = require('../repositories');

const queryBestPlayers = async (parent, args, context, info, schema) => {
  const { limit } = args;
  const key = `get_best_players_limit_${limit}`;
  const cachedData = await selectCachedQueryByKey({ key });

  if (cachedData) {
    const cachedDataTime = dayjs(cachedData.updated_at).utc();
    const daysSinceCachedDataUpdatedAt = dayjs.utc().diff(cachedDataTime, 'day');
    if (daysSinceCachedDataUpdatedAt < 1) {
      return cachedData.data;
    }
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

  await upsertCachedQueryData({ key, data: JSON.stringify(result) });

  return result;
};

module.exports = queryBestPlayers;
