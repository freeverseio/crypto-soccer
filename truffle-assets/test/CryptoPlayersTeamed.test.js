require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersTeamed');

contract('CryptoPlayersTeamed', (accounts) => {
    it('deployment', async () => {
        await CryptoPlayers.new().should.be.fulfilled;
    });

    // it('mint', async () => {
    //     const contract = await CryptoPlayers.new().should.be.fulfilled;
    //     const teamId = 1;
    //     const teamName = "panzerotto";
    //     await contract.mint(accounts[0], teamId);
    //     const name = await contract.getName(teamId).should.be.fulfilled;
    //     name.should.be.equal(name);
    // })
});
