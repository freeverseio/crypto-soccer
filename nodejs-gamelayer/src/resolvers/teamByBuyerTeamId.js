const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByBuyerTeamId = async (parent, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'teamByTeamId',
    args: {
      teamId: parent.buyerTeamId,
    },
    context,
    info,
  });
  if (!result) {
    return;
  }
  const teamName = await selectTeamName({ teamId: parent.buyerTeamId });
  const teamManagerName = await selectTeamManagerName({
    teamId: parent.buyerTeamId,
  });

  result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
  result.managerName =
    teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;

  return result;
};

module.exports = teamByBuyerTeamId;
