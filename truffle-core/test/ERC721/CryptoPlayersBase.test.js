require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersBaseMock');

contract('CryptoPlayersBase', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    });

    it('get name', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("");
    });
});
