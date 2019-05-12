require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('Storage');

contract('Storage', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await Storage.new().should.be.fulfilled;
    });

    it('initial number of team', async () => {
        const count = await instance.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('get name of invalid team', async () => {
        await instance.getTeamName(0).should.be.rejected;
    });

    it('get name of unexistent team', async () => {
        await instance.getTeamName(1).should.be.rejected;
    });
});