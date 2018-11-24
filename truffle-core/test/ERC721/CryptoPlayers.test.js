require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMock');

contract('CryptoPlayers', (accounts) => {
    it('deployment', async () => {
        await CryptoPlayers.new().should.be.fulfilled;
    });

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new();
        await contract.name().should.eventually.equal("CryptoSoccerPlayers");
        await contract.symbol().should.eventually.equal("CSP");
    });

    it('get state', async () => {
        const contract = await CryptoPlayers.new();
        const tokenId = 1;
        const state = 999;
        await contract.mint(accounts[0], tokenId, state).should.be.fulfilled;
        const result = await contract.getState(tokenId).should.be.fulfilled;
        result.toNumber().should.be.equal(state);
    });
});
