require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    })

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new();
        await contract.name().should.eventually.equal("CryptoSoccerPlayers");
        await contract.symbol().should.eventually.equal("CSP");
    });

    it('no initial players', async () => {
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });
});
