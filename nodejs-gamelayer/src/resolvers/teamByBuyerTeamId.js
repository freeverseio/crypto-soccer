const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByBuyerTeamId = (parent, args, context, info, schema) => {
  return info.mergeInfo
    .delegateToSchema({
      schema,
      operation: 'query',
      fieldName: 'teamByTeamId',
      args: {
        teamId: parent.buyerTeamId,
      },
      context,
      info,
    })
    .then((result) => {
      if (!result) {
        return;
      }
      const teamName = selectTeamName({ teamId: parent.buyerTeamId });
      const teamManagerName = selectTeamManagerName({
        teamId: parent.buyerTeamId,
      });

      result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
      result.managerName =
        teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;
      return result;
    });
};

module.exports = teamByBuyerTeamId;
