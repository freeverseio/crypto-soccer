require('chai')
    .use(require('chai-as-promised'))
    .should();

const TeamFactory = artifacts.require('TeamFactory');
const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    let contract = null;
    let teamFactory = null;

    beforeEach(async () => {
        teamFactory = await TeamFactory.new().should.be.fulfilled;
        contract = await CryptoTeams.new(teamFactory.address).should.be.fulfilled;
    });

    it('TeamFactory address', async () => {
        const result = await contract.getTeamFactory().should.be.fulfilled;
        result.should.be.equal(teamFactory.address);
    }); 
});
