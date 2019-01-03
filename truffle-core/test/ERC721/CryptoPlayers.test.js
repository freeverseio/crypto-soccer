require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    })
});
