const { selectTeamName, selectTeamManagerName } = require('../repositories');

const queryTeamByTeamId = async (parent, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'teamByTeamId',
    args: {
      teamId: args.teamId,
    },
    context,
    info,
  });

  if (!result) {
    return;
  }
  const teamName = await selectTeamName({ teamId: args.teamId });
  const teamManagerName = await selectTeamManagerName({
    teamId: args.teamId,
  });

  result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
  result.managerName =
    teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;

  return result;
};

module.exports = queryTeamByTeamId;
