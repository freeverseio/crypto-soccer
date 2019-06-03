module.exports = function Resolvers(universe) {
  this.Query = {
    countTeams: () => universe.countTeams(),
    allTeams: () => universe.getTeamIds(),
    getTeam: (_, { id }) => id,
    getPlayer: (_, { id }) => id
  };

  this.Team = {
    id: (id) => id,
    name: (id) => universe.getTeamName(id),
    players: (id) => universe.getTeamPlayerIds(id)
  };

  this.Player = {
    id: (id) => id,
    name: (id) => universe.getPlayerName(id),
    defence: (id) => universe.getPlayerDefence(id),
    speed: (id) => universe.getPlayerSpeed(id),
    pass: (id) => universe.getPlayerPass(id),
    shoot: (id) => universe.getPlayerShoot(id),
    endurance: (id) => universe.getPlayerEndurance(id),
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