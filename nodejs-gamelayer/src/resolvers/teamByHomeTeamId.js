const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByHomeTeamId = (match, args, context, info, schema) => {
  return info.mergeInfo
    .delegateToSchema({
      schema,
      operation: 'query',
      fieldName: 'teamByTeamId',
      args: {
        teamId: match.homeTeamId,
      },
      context,
      info,
    })
    .then((result) => {
      if (!result) {
        return;
      }
      const teamName = selectTeamName({ teamId: match.homeTeamId });
      const teamManagerName = selectTeamManagerName({
        teamId: match.homeTeamId,
      });

      result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
      result.managerName =
        teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;
      return result;
    });
};

module.exports = teamByHomeTeamId;
