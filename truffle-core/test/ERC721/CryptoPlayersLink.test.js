require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersLink');

contract('CryptoPlayersLink', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    });

    it('set team not possible by any account', async () =>{
        const playerId = 1;
        const teamId = 1;
        await contract.mintWithName(accounts[0], playerId, "player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.rejected;
    });
});
