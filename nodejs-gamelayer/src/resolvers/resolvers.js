const Web3 = require('web3');
const { selectTeamName, selectTeamManagerName, updateTeamName, updateTeamManagerName } = require('../repositories');
const { TeamValidation, BidValidation } = require('../validations');
const getMessagesResolver = require('./getMessagesResolver');
const setMessageReadResolver = require('./setMessageReadResolver');
const setMailboxStartResolver = require('./setMailboxStartResolver');
const setBroadcastMessageResolver = require('./setBroadcastMessageResolver');
const setMessageResolver = require('./setMessageResolver');
const teamByTeamId = require('./teamByTeamId');
const teamByHomeTeamId = require('./teamByHomeTeamId');
const teamByVisitorTeamId = require('./teamByVisitorTeamId');
const teamByBuyerTeamId = require('./teamByBuyerTeamId');
const getNumUnreadMessagesResolver = require('./getNumUnreadMessagesResolver');
const getLastTimeLoggedInResolver = require('./getLastTimeLoggedIn');
const setLastTimeLoggedInResolver = require('./setLastTimeLoggedIn');

const web3 = new Web3('');

const resolvers = ({ horizonRemoteSchema }) => {
  return {
    TeamsHistory: {
      name: {
        fragment: `... on TeamsHistory { teamId }`,
        resolve(team, args, context, info) {
          return selectTeamName({ teamId: team.teamId }).then((result) => {
            return result && result.team_name ? result.team_name : team.name;
          });
        },
      },
    },
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
      createBid: async (_, args, context, info) => {
        const {
          input: { teamId, rnd, auctionId, extraPrice, signature },
        } = args;
        const bidValidation = new BidValidation({ teamId, rnd, auctionId, extraPrice, signature, web3 });
        const isAllowed = await bidValidation.isAllowedToBid();

        if (!isAllowed) {
          return 'User not allowed to bid for that amount';
        } else {
          return info.mergeInfo.delegateToSchema({
            schema: horizonRemoteSchema,
            operation: 'mutation',
            fieldName: 'createBid',
            args,
            context,
            info,
          });
        }
      },
      setMessage: setMessageResolver,
      setBroadcastMessage: setBroadcastMessageResolver,
      setMailboxStart: setMailboxStartResolver,
      setMessageRead: setMessageReadResolver,
      setLastTimeLoggedIn: setLastTimeLoggedInResolver,
    },
    Query: {
      getMessages: getMessagesResolver,
      getNumUnreadMessages: getNumUnreadMessagesResolver,
      getLastTimeLoggedIn: getLastTimeLoggedInResolver,
    },
  };
};

module.exports = resolvers;
