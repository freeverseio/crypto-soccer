require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    const name = "name";
    const symbol = "symbol";

    it('deployment', async () => {
        await CryptoTeams.new(name, symbol).should.be.fulfilled;
    });

    it('mint a team', async () => {

    });
});
