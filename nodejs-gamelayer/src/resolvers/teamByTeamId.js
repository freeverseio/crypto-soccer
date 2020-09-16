const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByTeamId = (parent, args, context, info, schema) => {
  return info.mergeInfo
    .delegateToSchema({
      schema,
      operation: 'query',
      fieldName: 'teamByTeamId',
      args: {
        teamId: parent.teamId,
      },
      context,
      info,
    })
    .then(async (result) => {
      if (!result) {
        return;
      }
      const teamName = await selectTeamName({ teamId: result.teamId });
      const teamManagerName = await selectTeamManagerName({
        teamId: result.teamId,
      });
      console.log(`I have teamName from horizon: ${result.teamName} and temName from GameLayer: ${teamName}`);
      console.log(
        `I have teamManagerName from horizon: ${result.teamManagerName} and temName from GameLayer: ${teamManagerName}`
      );
      result.teamName = teamName ? teamName : result.teamName;
      result.teamManagerName = teamManagerName ? teamManagerName : result.teamManagerName;
      return result;
    });
};

module.exports = teamByTeamId;
