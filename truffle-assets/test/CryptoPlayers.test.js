require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    it('deployment', async () => {
        const name = "name";
        const symbol = "symbol";
        await CryptoPlayers.new(name, symbol).should.be.fulfilled;
    });
});
