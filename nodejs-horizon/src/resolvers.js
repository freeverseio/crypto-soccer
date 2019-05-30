module.exports = function Resolvers(universe) {
  this.Query = {
    countTeams: () => universe.countTeams(),
    allTeams: (parent, args, context, info) => [{ id: 3 }],
    getTeam: (_, { id }) => id,
    getPlayer: (_, { id }) => id
  };

  this.Team = {
    name: (id) => universe.getTeamName(id),
    players: (id) => universe.getTeamPlayerIds(id)
  };

  this.Player = {
    id: (id) => id,
    name: (id) => universe.getPlayerName(id)
  }

  this.Mutation = {
    createTeam: (parent, args, context, info) => universe.createTeam(args.name, args.owner)
  };

  this.Subscription = {
    // teamCreated: {
    //   subscribe: () => pubsub.asyncIterator([TEAM_CREATED])
    // }
  };
}