require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });
});
