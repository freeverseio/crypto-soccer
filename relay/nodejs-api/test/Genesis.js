const Web3 = require('web3');
const statesJSON = require('../contracts/TeamState.json');
const assetsJSON = require('../contracts/Assets.json');
const engineJSON = require('../contracts/Engine.json');
const leaguesJSON = require('../contracts/Leagues.json');

const Universe = async (provider, from) => {
    web3 = new Web3(provider, null, {});

    const States = new web3.eth.Contract(statesJSON.abi);
    let gas = await States.deploy({ data: statesJSON.bytecode }).estimateGas();
    states = await States.deploy({ data: statesJSON.bytecode }).send({ from, gas });

    const Assets = new web3.eth.Contract(assetsJSON.abi);
    gas = await Assets.deploy({ data: assetsJSON.bytecode, arguments: [states.options.address] }).estimateGas();
    assets = await Assets.deploy({ data: assetsJSON.bytecode, arguments: [states.options.address] }).send({ from, gas });

    const Engine = new web3.eth.Contract(engineJSON.abi);
    gas = await Engine.deploy({ data: engineJSON.bytecode }).estimateGas();
    engine = await Engine.deploy({ data: engineJSON.bytecode }).send({ from, gas });

    const Leagues = new web3.eth.Contract(leaguesJSON.abi);
    gas = await Leagues.deploy({ data: leaguesJSON.bytecode, arguments: [states.options.address, engine.options.address] }).estimateGas();
    leagues = await Leagues.deploy({ data: leaguesJSON.bytecode, arguments: [states.options.address, engine.options.address] }).send({ from, gas });

    return {
        states,
        assets,
        leagues
    }
}

module.exports = Universe;