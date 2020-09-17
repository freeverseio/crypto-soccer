const { selectTeamName, selectTeamManagerName } = require('../repositories');

const teamByTeamId = async (parent, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'teamByTeamId',
    args: {
      teamId: parent.teamId,
    },
    context,
    info,
  });

  if (!result) {
    return;
  }

  const teamName = await selectTeamName({ teamId: parent.teamId });
  const teamManagerName = await selectTeamManagerName({
    teamId: parent.teamId,
  });

  result.name = teamName && teamName.team_name ? teamName.team_name : result.name;
  result.managerName =
    teamManagerName && teamManagerName.team_manager_name ? teamManagerName.team_manager_name : result.managerName;
  return result;
};
// const teamByTeamId = (parent, args, context, info, schema) => {
//   return info.mergeInfo
//     .delegateToSchema({
//       schema,
//       operation: 'query',
//       fieldName: 'teamByTeamId',
//       args: {
//         teamId: parent.teamId,
//       },
//       context,
//       info,
//     })
//     .then(async (result) => {
//       if (!result) {
//         return;
//       }
//       const teamName = await selectTeamName({ teamId: result.teamId });
//       const teamManagerName = await selectTeamManagerName({
//         teamId: result.teamId,
//       });
//       console.log(`I have teamName from horizon: ${result.teamName} and temName from GameLayer: ${teamName}`);
//       console.log(
//         `I have teamManagerName from horizon: ${result.teamManagerName} and temName from GameLayer: ${teamManagerName}`
//       );
//       result.teamName = teamName ? teamName : result.teamName;
//       result.teamManagerName = teamManagerName ? teamManagerName : result.teamManagerName;
//       return result;
//     });
// };

module.exports = teamByTeamId;
