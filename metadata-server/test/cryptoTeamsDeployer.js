const Web3 = require('web3');
const contractJSON = require('../../truffle-core/build/contracts/CryptoTeams.json')

module.exports = async ({ provider, playersContract, sender }) => {
    const web3 = new Web3(provider);
    const contract = new web3.eth.Contract(contractJSON.abi);

    const instance = await contract.deploy({
        data: contractJSON.bytecode,
        arguments: [playersContract.options.address]
    })
        .send({
            from: sender,
            gas: 4712388,
            gasPrice: provider.gasPrice
        })
        .on('error', console.error)
        .catch(console.error);

    await playersContract.methods.addTeamsContract(instance.options.address).send({
        from: sender,
        gas: 4712388,
        gasPrice: provider.gasPrice
    });

    return instance;
}
