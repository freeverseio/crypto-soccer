module.exports = function Resolvers(universe) {
  this.Query = {
    countTeams: () => universe.assets.methods.countTeams().call(),
    allTeams: async () => {
      const count = await universe.assets.methods.countTeams().call();
      let ids = [];
      for (let i = 1; i <= count; i++)
        ids.push(i);
      return ids;
    },
    getTeam: (_, { id }) => id,
    getPlayer: (_, { id }) => id,
    countLeagues: () => universe.leagues.methods.leaguesCount().call(),
  };

  this.Mutation = {
    createTeam: async (_, { name, owner }) => {
      const gas = await universe.assets.methods.createTeam(name, owner).estimateGas();
      await universe.assets.methods.createTeam(name, owner).send({ from: universe.from, gas });
    },
    createLeague: async (_, { initBlock, step, teamIds, tactics }) => {
      const { leagues, from } = universe;
      const count = await leagues.methods.leaguesCount().call();
      const id = count + 1;
      const gas = await leagues.methods.create(
        id,
        initBlock,
        step,
        teamIds,
        tactics
      ).estimateGas();
      await leagues.methods.create(
        id,
        initBlock,
        step,
        teamIds,
        tactics
      ).send({ from, gas });
    },
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
    name: (id) => "player_" + id,
    defence: async (id) => {
      const state = await universe.assets.methods.getPlayerState(id).call();
      return await universe.playerState.methods.getDefence(state).call();
    },
    speed: async (id) => {
      const state = await universe.assets.methods.getPlayerState(id).call();
      return await universe.playerState.methods.getSpeed(state).call();
    },
    pass: async (id) => {
      const state = await universe.assets.methods.getPlayerState(id).call();
      return await universe.playerState.methods.getPass(state).call();
    },
    shoot: async (id) => {
      const state = await universe.assets.methods.getPlayerState(id).call();
      return await universe.playerState.methods.getShoot(state).call();
    },
    endurance: async (id) => {
      const state = await universe.assets.methods.getPlayerState(id).call();
      return await universe.playerState.methods.getEndurance(state).call();
    },
    team: async (id) => {
      const state = await universe.assets.methods.getPlayerState(id).call();
      return await universe.playerState.methods.getCurrentTeamId(state).call();
    },
  };

  this.League = {
    id: (id) => id,
    initBlock: (id) => universe.legues.methods.getInitBlock(id).call(),
    step: (id) => universe.leagues.methods.getStep(id).call(),
    nTeams: (id) => universe.leagues.methods.getNTeams(id).call(),
  };
}