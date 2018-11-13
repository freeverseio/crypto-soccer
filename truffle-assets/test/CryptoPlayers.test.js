require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    it('deployment', async () => {
        await CryptoPlayers.new().should.be.fulfilled;
    });
});
