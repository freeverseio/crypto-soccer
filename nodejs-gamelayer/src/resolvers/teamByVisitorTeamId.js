const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByVisitorTeamId = (match, args, context, info, schema) => {
  return info.mergeInfo
    .delegateToSchema({
      schema,
      operation: 'query',
      fieldName: 'teamByTeamId',
      args: {
        teamId: match.visitorTeamId,
      },
      context,
      info,
    })
    .then((result) => {
      if (!result) {
        return;
      }
      const teamName = selectTeamName({ teamId: match.visitorTeamId });
      const teamManagerName = selectTeamManagerName({
        teamId: match.visitorTeamId,
      });

      result.name = teamName && teamName.team_name ? teamName : result.name;
      result.managerName =
        teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;
      return result;
    });
};

module.exports = teamByVisitorTeamId;
