const Web3 = require('web3');
const contractJSON = require('../../truffle-core/build/contracts/Players.json')

module.exports = async ({ provider, sender }) => {
    const web3 = new Web3(provider);
    const contract = new web3.eth.Contract(contractJSON.abi);

    const instance = await contract.deploy({
        data: contractJSON.bytecode
    })
        .send({
            from: sender,
            gas: 4712388,
            gasPrice: provider.gasPrice
        })
        .on('error', console.error)
        .catch(console.error);

    return instance;
}
