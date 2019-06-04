module.exports = function Resolvers(universe) {
  this.Query = {
    countTeams: () => universe.countTeams(),
    allTeams: () => universe.getTeamIds(),
    getTeam: (_, { id }) => id,
    getPlayer: (_, { id }) => id,
    countLeagues: () => universe.leagues.methods.leaguesCount().call(),
  };

  this.Mutation = {
    createTeam: (_, { name, owner }) => universe.createTeam(name, owner),
    createLeague: (_, { initBlock, step, teamIds, tactics }) => universe.createLeague(initBlock, step, teamIds, tactics),
  };

  this.Subscription = {
    // teamCreated: {
    //   subscribe: () => pubsub.asyncIterator([TEAM_CREATED])
    // }
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
    team: (id) => universe.getPlayerTeamId(id),
  };

  this.League = {
    id: (id) => id,
    initBlock: (id) => universe.legues.methods.getInitBlock(id).call(),
    step: (id) => universe.leagues.methods.getStep(id).call(),
    nTeams: (id) => universe.leagues.methods.getNTeams(id).call(),
  };
}