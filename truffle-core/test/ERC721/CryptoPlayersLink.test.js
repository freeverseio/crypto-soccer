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
        const teamId = 1;
        await contract.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.rejected;
    });

    it('set team possible by team account', async () => {
        const teamId = 1;
        await contract.setTeamsContract(accounts[0]).should.be.fulfilled;
        await contract.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.fulfilled;
    })

    it('get team contract', async () => {
        const address = "0x0000000000000000000000000000000000000001";
        await contract.setTeamsContract(address).should.be.fulfilled;
        const result = await contract.getTeamsContract().should.be.fulfilled;
        result.should.be.equal(address);
    });
});
