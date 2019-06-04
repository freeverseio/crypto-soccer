const Web3 = require('web3');
const playerStateJSON = require('../../truffle-core/build/contracts/PlayerState.json');
const assetsJSON = require('../../truffle-core/build/contracts/Assets.json');
const engineJSON = require('../../truffle-core/build/contracts/Engine.json');
const leaguesJSON = require('../../truffle-core/build/contracts/Leagues.json');

function Resolvers({
  provider,
  playerStateAddress,
  assetsAddress,
  leaguesAddress,
  from
}) {
  this.web3 = new Web3(provider, null, {});
  this.playerState = new this.web3.eth.Contract(playerStateJSON.abi, playerStateAddress);
  this.assets = new this.web3.eth.Contract(assetsJSON.abi, assetsAddress);
  this.leagues = new this.web3.eth.Contract(leaguesJSON.abi, leaguesAddress);
  this.from = from;

  this.Query = {
    countTeams: () => this.assets.methods.countTeams().call(),
    allTeams: async () => {
      const count = await this.assets.methods.countTeams().call();
      let ids = [];
      for (let i = 1; i <= count; i++)
        ids.push(i);
      return ids;
    },
    getTeam: (_, { id }) => id,
    getPlayer: (_, { id }) => id,
    countLeagues: () => this.leagues.methods.leaguesCount().call(),
  };

  this.Mutation = {
    createTeam: async (_, { name, owner }) => {
      const gas = await this.assets.methods.createTeam(name, owner).estimateGas();
      await this.assets.methods.createTeam(name, owner).send({ from: this.from, gas });
    },
    createLeague: async (_, { initBlock, step, teamIds, tactics }) => {
      const { leagues, from } = this;
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
    name: (id) => this.getTeamName(id),
    players: (id) => this.getTeamPlayerIds(id)
  };

  this.Player = {
    id: (id) => id,
    name: (id) => "player_" + id,
    defence: async (id) => {
      const state = await this.assets.methods.getPlayerState(id).call();
      return await this.playerState.methods.getDefence(state).call();
    },
    speed: async (id) => {
      const state = await this.assets.methods.getPlayerState(id).call();
      return await this.playerState.methods.getSpeed(state).call();
    },
    pass: async (id) => {
      const state = await this.assets.methods.getPlayerState(id).call();
      return await this.playerState.methods.getPass(state).call();
    },
    shoot: async (id) => {
      const state = await this.assets.methods.getPlayerState(id).call();
      return await this.playerState.methods.getShoot(state).call();
    },
    endurance: async (id) => {
      const state = await this.assets.methods.getPlayerState(id).call();
      return await this.playerState.methods.getEndurance(state).call();
    },
    team: async (id) => {
      const state = await this.assets.methods.getPlayerState(id).call();
      return await this.playerState.methods.getCurrentTeamId(state).call();
    },
  };

  this.League = {
    id: (id) => id,
    initBlock: (id) => this.legues.methods.getInitBlock(id).call(),
    step: (id) => this.leagues.methods.getStep(id).call(),
    nTeams: (id) => this.leagues.methods.getNTeams(id).call(),
  };
}

module.exports = Resolvers;