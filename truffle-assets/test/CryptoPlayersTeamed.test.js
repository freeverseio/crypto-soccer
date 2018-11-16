require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersTeamed');

contract('CryptoPlayersTeamed', (accounts) => {
    it('deployment', async () => {
        await CryptoPlayers.new().should.be.fulfilled;
    });
});
