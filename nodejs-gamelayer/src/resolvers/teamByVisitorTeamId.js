const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByVisitorTeamId = async (match, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'teamByTeamId',
    args: {
      teamId: match.visitorTeamId,
    },
    context,
    info,
  });
  if (!result) {
    return;
  }
  const teamName = await selectTeamName({ teamId: match.visitorTeamId });
  const teamManagerName = await selectTeamManagerName({
    teamId: match.visitorTeamId,
  });

  result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
  result.managerName =
    teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;

  return result;
};

module.exports = teamByVisitorTeamId;
