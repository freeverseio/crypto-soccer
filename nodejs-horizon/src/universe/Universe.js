const Web3 = require('web3');
const playerStateJSON = require('../../../truffle-core/build/contracts/PlayerState.json');
const assetsJSON = require('../../../truffle-core/build/contracts/Assets.json');

class Universe {
    constructor(provider, assetsAddress, from) {
        this.web3 = new Web3(provider, null, {});
        this.assets = new this.web3.eth.Contract(assetsJSON.abi, assetsAddress);
        this.from = from;
    }

    async genesis() {
        const { web3, from } = this;

        const PlayerState = new web3.eth.Contract(playerStateJSON.abi);
        let gas = await PlayerState.deploy({ data: playerStateJSON.bytecode }).estimateGas();
        const playerState = await PlayerState.deploy({ data: playerStateJSON.bytecode }).send({ from, gas });

        const Assets = new web3.eth.Contract(assetsJSON.abi);
        gas = await Assets.deploy({ data: assetsJSON.bytecode, arguments: [playerState.options.address] }).estimateGas();
        const assets = await Assets.deploy({ data: assetsJSON.bytecode, arguments: [playerState.options.address] }).send({ from, gas });

        this.playerState = playerState;
        this.assets = assets;
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
}

module.exports = Universe;