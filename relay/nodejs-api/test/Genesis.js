const Web3 = require('web3');
const assetsJSON = require('../contracts/Assets.json');
const marketJSON = require('../contracts/Market.json');

const Universe = async (provider, from) => {
    web3 = new Web3(provider, null, {});

    const Assets = new web3.eth.Contract(assetsJSON.abi);
    let gas = await Assets.deploy({ data: assetsJSON.bytecode }).estimateGas();
    assets = await Assets.deploy({ data: assetsJSON.bytecode }).send({ from, gas });
    await assets.methods.init().send({ from, gas });

    const Market = new web3.eth.Contract(marketJSON.abi);
    gas = await Market.deploy({ data: marketJSON.bytecode }).estimateGas();
    market = await Market.deploy({ data: marketJSON.bytecode }).send({ from, gas });

    gas = await market.methods.setAssetsAddress(assets.options.address).estimateGas();
    await market.methods.setAssetsAddress(assets.options.address).send({ from, gas });

    return {
        assets,
        market
    }
}

module.exports = Universe;