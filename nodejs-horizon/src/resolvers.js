function Resolvers({
  states,
  assets,
  leagues,
  from
}) {
  this.Query = {
    countTeams: () => assets.methods.countTeams().call(),
    allTeams: async () => {
      const count = await assets.methods.countTeams().call();
      let ids = [];
      for (let i = 1; i <= count; i++)
        ids.push(i);
      return ids;
    },
    getTeam: (_, { id }) => id,
    getPlayer: (_, { id }) => id,
    countLeagues: () => leagues.methods.leaguesCount().call(),
  };

  this.Mutation = {
    createTeam: async (_, { name, owner }) => {
      const gas = await assets.methods.createTeam(name, owner).estimateGas();
      const receipt = await assets.methods.createTeam(name, owner).send({ from: from, gas });
      return receipt.events.TeamCreation.returnValues.teamId;
    },
    createLeague: async (_, { initBlock, step, teamIds, tactics }) => {
      const count = await leagues.methods.leaguesCount().call();
      const id = count;
      const gas = await leagues.methods.create(
        id,
        initBlock,
        step,
        teamIds,
        tactics
      ).estimateGas();
      const receipt = await leagues.methods.create(
        id,
        initBlock,
        step,
        teamIds,
        tactics
      ).send({ from, gas });
      return receipt.events.LeagueCreated.returnValues.id;
    },
  };

  this.Subscription = {
    // teamCreated: {
    //   subscribe: () => pubsub.asyncIterator([TEAM_CREATED])
    // }
  };

  this.Team = {
    id: (id) => id,
    name: (id) => assets.methods.getTeamName(id).call(),
    players: (id) => assets.methods.getTeamPlayerIds(id).call()
  };

  this.Player = {
    id: (id) => id,
    name: (id) => "player_" + id,
    defence: async (id) => {
      const state = await assets.methods.getPlayerState(id).call();
      return await states.methods.getDefence(state).call();
    },
    speed: async (id) => {
      const state = await assets.methods.getPlayerState(id).call();
      return await states.methods.getSpeed(state).call();
    },
    pass: async (id) => {
      const state = await assets.methods.getPlayerState(id).call();
      return await states.methods.getPass(state).call();
    },
    shoot: async (id) => {
      const state = await assets.methods.getPlayerState(id).call();
      return await states.methods.getShoot(state).call();
    },
    endurance: async (id) => {
      const state = await assets.methods.getPlayerState(id).call();
      return await states.methods.getEndurance(state).call();
    },
    team: async (id) => {
      const state = await assets.methods.getPlayerState(id).call();
      return await states.methods.getCurrentTeamId(state).call();
    },
  };

  this.League = {
    id: (id) => id,
    initBlock: (id) => legues.methods.getInitBlock(id).call(),
    step: (id) => leagues.methods.getStep(id).call(),
    nTeams: (id) => leagues.methods.getNTeams(id).call(),
  };
}

module.exports = Resolvers;