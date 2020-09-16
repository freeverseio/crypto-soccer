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
      const teamName = selectTeamName({ teamId: result.teamId });
      const teamManagerName = selectTeamManagerName({
        teamId: result.teamId,
      });

      result.teamName = teamName ? teamName : result.teamName;
      result.teamManagerName = teamManagerName ? teamManagerName : result.teamManagerName;
      return result;
    });
};

module.exports = queryTeamByTeamId;
