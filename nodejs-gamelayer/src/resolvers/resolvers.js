const Web3 = require('web3');
const { selectTeamName, selectTeamManagerName, updateTeamName, updateTeamManagerName } = require('../repositories');
const { TeamValidation } = require('../validations');

const web3 = new Web3('');

const resolvers = ({ horizonRemoteSchema }) => {
  return {
    Team: {
      name: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return selectTeamName({ teamId: team.teamId }).then((result) => {
            return result && result.team_name ? result.team_name : team.name;
          });
        },
      },
      managerName: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return selectTeamManagerName({ teamId: team.teamId }).then((result) => {
            return result && result.team_manager_name ? result.team_manager_name : team.managerName;
          });
        },
      },
    },
    Match: {
      teamByHomeTeamId: {
        fragment: `... on Match { homeTeamId }`,
        resolve(match, args, context, info) {
          console.log('Intercepted!', match);
          console.log('resolve -> match.homeTeamId', match.homeTeamId);
          return info.mergeInfo
            .delegateToSchema({
              schema: horizonRemoteSchema,
              operation: 'query',
              fieldName: 'teamByTeamId',
              args: {
                teamId: match.homeTeamId,
              },
              context,
              info,
            })
            .then((result) => {
              console.log('hola=???', result);
              const teamName = selectTeamName({ teamId: result.teamId });
              const teamManagerName = selectTeamManagerName({
                teamId: result.teamId,
              });

              result.teamName = teamName ? teamName : result.teamName;
              result.teamManagerName = teamManagerName ? teamManagerName : result.teamManagerName;
              return result;
            });
        },
      },
    },
    Mutation: {
      setTeamName: async (_, { input: { teamId, name, signature } }) => {
        const teamValidation = new TeamValidation({
          teamId,
          name,
          signature,
          web3,
        });
        const isSignerOwner = await teamValidation.isSignerOwner();

        if (isSignerOwner) {
          await updateTeamName({ teamId, teamName: name });
          return teamId;
        } else {
          return 'Signer is not the team owner';
        }
      },
      setTeamManagerName: async (_, { input: { teamId, name, signature } }) => {
        const teamValidation = new TeamValidation({
          teamId,
          name,
          signature,
          web3,
        });
        const isSignerOwner = await teamValidation.isSignerOwner();

        if (isSignerOwner) {
          await updateTeamManagerName({ teamId, teamManagerName: name });
          return teamId;
        } else {
          return 'Signer is not the team owner';
        }
      },
    },
  };
};

module.exports = resolvers;
