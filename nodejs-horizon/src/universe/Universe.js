const Web3 = require('web3');
const playerStateJSON = require('../../../truffle-core/build/contracts/PlayerState.json');
const assetsJSON = require('../../../truffle-core/build/contracts/Assets.json');
const engineJSON = require('../../../truffle-core/build/contracts/Engine.json');
const leaguesJSON = require('../../../truffle-core/build/contracts/Leagues.json');

class Universe {
    constructor({
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
    }

    async genesis() {
        const { web3, from } = this;

        const PlayerState = new web3.eth.Contract(playerStateJSON.abi);
        let gas = await PlayerState.deploy({ data: playerStateJSON.bytecode }).estimateGas();
        this.playerState = await PlayerState.deploy({ data: playerStateJSON.bytecode }).send({ from, gas });

        const Assets = new web3.eth.Contract(assetsJSON.abi);
        gas = await Assets.deploy({ data: assetsJSON.bytecode, arguments: [this.playerState.options.address] }).estimateGas();
        this.assets = await Assets.deploy({ data: assetsJSON.bytecode, arguments: [this.playerState.options.address] }).send({ from, gas });

        const Engine = new web3.eth.Contract(engineJSON.abi);
        gas = await Engine.deploy({ data: engineJSON.bytecode }).estimateGas();
        this.engine = await Engine.deploy({ data: engineJSON.bytecode }).send({ from, gas });

        const Leagues = new web3.eth.Contract(leaguesJSON.abi);
        gas = await Leagues.deploy({ data: leaguesJSON.bytecode, arguments: [this.playerState.options.address, this.engine.options.address] }).estimateGas();
        this.leagues = await Leagues.deploy({ data: leaguesJSON.bytecode, arguments: [this.playerState.options.address, this.engine.options.address] }).send({ from, gas });
    }

    async countTeams() {
        return await this.assets.methods.countTeams().call();
    }

    async createTeam(name, owner) {
        const { assets, from } = this;
        const gas = await assets.methods.createTeam(name, owner).estimateGas();
        await assets.methods.createTeam(name, owner).send({ from, gas });
    }

    async getTeamName(id) {
        return await this.assets.methods.getTeamName(id).call();
    }

    async getTeamPlayerIds(id) {
        return await this.assets.methods.getTeamPlayerIds(id).call();
    }

    getPlayerName(id) {
        return "player_" + id;
    }

    async getPlayerDefence(id) {
        const state = await this.assets.methods.getPlayerState(id).call();
        return await this.playerState.methods.getDefence(state).call();
    }

    async getPlayerSpeed(id) {
        const state = await this.assets.methods.getPlayerState(id).call();
        return await this.playerState.methods.getSpeed(state).call();
    }

    async getPlayerPass(id) {
        const state = await this.assets.methods.getPlayerState(id).call();
        return await this.playerState.methods.getPass(state).call();
    }

    async getPlayerShoot(id) {
        const state = await this.assets.methods.getPlayerState(id).call();
        return await this.playerState.methods.getShoot(state).call();
    }

    async getPlayerEndurance(id) {
        const state = await this.assets.methods.getPlayerState(id).call();
        return await this.playerState.methods.getEndurance(state).call();
    }

    async getTeamIds() {
        const count = await this.assets.methods.countTeams().call();
        let ids = [];
        for (let i = 1 ; i <= count ; i++)
            ids.push(i);
        return ids;
    }

    async getPlayerTeamId(id) {
        const state = await this.assets.methods.getPlayerState(id).call();
        return await this.playerState.methods.getCurrentTeamId(state).call();
    }

    async createLeague(initBlock, step, teamIds, tactics) {
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
    }
}

module.exports = Universe;