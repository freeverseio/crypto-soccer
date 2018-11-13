require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    it('deployment', async () => {
        await CryptoPlayers.new({gas: 6500000}).should.be.fulfilled;
    });
});
