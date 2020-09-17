const { selectTeamName, selectTeamManagerName } = require('../repositories');

const queryTeamByTeamId = (parent, args, context, info, schema) => {
  return info.mergeInfo
    .delegateToSchema({
      schema,
      operation: 'query',
      fieldName: 'teamByTeamId',
      args: {
        teamId: args.teamId,
      },
      context,
      info,
    })
    .then((result) => {
      if (!result) {
        return;
      }
      const teamName = selectTeamName({ teamId: args.teamId });
      const teamManagerName = selectTeamManagerName({
        teamId: args.teamId,
      });

      result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
      result.managerName =
        teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;
      return result;
    });
};

module.exports = queryTeamByTeamId;
