const Web3 = require('web3');
const marketJSON = require('../contracts/Market.json');

const Universe = async (provider, from) => {
    web3 = new Web3(provider, null, {});

    const Market = new web3.eth.Contract(marketJSON.abi);
    let gas = await Market.deploy({ data: marketJSON.bytecode }).estimateGas();
    market = await Market.deploy({ data: marketJSON.bytecode }).send({ from, gas });
    await market.methods.init().send({ from, gas });

    return {
        market
    }
}

module.exports = Universe;