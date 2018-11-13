require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    const name = "name";
    const symbol = "symbol";

    it('deployment', async () => {
        await CryptoPlayers.new(name, symbol).should.be.fulfilled;
    });

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new(name, symbol);
        await contract.name().should.eventually.equal(name);
        await contract.symbol().should.eventually.equal(symbol);
    });
});
