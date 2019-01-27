const Web3 = require('web3');
const ganache = require("ganache-cli");
const cryptoPlayersDeployer = require('./cryptoPlayersDeployer');
const cryptoTeamsDeployer = require('./cryptoTeamsDeployer');

module.exports.deployer = async (identity) => {
    const provider = ganache.provider({
        accounts: [{
            secretKey: identity.privateKey,
            balance: Web3.utils.toWei('100', 'ether')
        }]
    });

    const playersContract = await cryptoPlayersDeployer({ provider, sender: identity.address });
    const teamsContract = await cryptoTeamsDeployer({ provider, playersContract, sender: identity.address });

    await playersContract.methods.mint(identity.address, "player").send({
        from: identity.address,
        gas: 4712388,
        gasPrice: provider.gasPrice
    }).should.be.fulfilled;
    const playerId = await playersContract.methods.getPlayerId("player").call().should.be.fulfilled;

    await teamsContract.methods.mint(identity.address, "team").send({
        from: identity.address,
        gas: 4712388,
        gasPrice: provider.gasPrice
    }).should.be.fulfilled;
    const teamId = await teamsContract.methods.getTeamId("team").call().should.be.fulfilled;

    await teamsContract.methods.addPlayer(teamId, playerId).send({
        from: identity.address,
        gas: 4712388,
        gasPrice: provider.gasPrice
    }).should.be.fulfilled;

    return {playersContract, teamsContract, playerId, teamId};
}
