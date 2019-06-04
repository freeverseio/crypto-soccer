const Web3 = require('web3');
const playerStateJSON = require('../../../truffle-core/build/contracts/PlayerState.json');
const assetsJSON = require('../../../truffle-core/build/contracts/Assets.json');
const engineJSON = require('../../../truffle-core/build/contracts/Engine.json');
const leaguesJSON = require('../../../truffle-core/build/contracts/Leagues.json');

const Universe = async (web3, from) => {
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

module.exports = Universe;