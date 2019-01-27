const Web3 = require('web3');
const ganache = require("ganache-cli");
const cryptoPlayersDeployer = require('./cryptoPlayersDeployer');
const cryptoTeamsDeployer = require('./cryptoTeamsDeployer');

let provider = null;
let playersContract = null;
let teamsContract = null;

module.exports.mintPlayer = async (sender, name) => {
    await playersContract.methods.mint(sender, name).send({
        from: sender,
        gas: 4712388,
        gasPrice: provider.gasPrice
    }).should.be.fulfilled;
    const playerId = await playersContract.methods.getPlayerId(name).call().should.be.fulfilled;
    return playerId;
}

module.exports.addPlayer = async (sender, teamId, playerId) => {
    await teamsContract.methods.addPlayer(teamId, playerId).send({
        from: sender,
        gas: 4712388,
        gasPrice: provider.gasPrice
    }).should.be.fulfilled;
}

module.exports.deployer = async (identity) => {
    provider = ganache.provider({
        accounts: [{
            secretKey: identity.privateKey,
            balance: Web3.utils.toWei('100', 'ether')
        }]
    });

    playersContract = await cryptoPlayersDeployer({ provider, sender: identity.address });
    teamsContract = await cryptoTeamsDeployer({ provider, playersContract, sender: identity.address });
    const playerId = await module.exports.mintPlayer(identity.address, "player").should.be.fulfilled;

    await teamsContract.methods.mint(identity.address, "team").send({
        from: identity.address,
        gas: 4712388,
        gasPrice: provider.gasPrice
    }).should.be.fulfilled;
    const teamId = await teamsContract.methods.getTeamId("team").call().should.be.fulfilled;

    module.exports.addPlayer(identity.address, teamId, playerId).should.be.fulfilled;

    return {playersContract, teamsContract, playerId, teamId};
}
