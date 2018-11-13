require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    it('deployment', async () => {
        await CryptoTeams.new().should.be.fulfilled;
    });
});
