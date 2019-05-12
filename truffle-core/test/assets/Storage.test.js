require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('StorageMock');

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

    it('existence of null player', async () => {
        const exists = await instance.playerExists(playerId = 0).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('existence of unexistent player', async () => {
        const exists = await instance.playerExists(playerId = 1).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('existence of existent player', async () => {
        await instance.addTeam("Barca").should.be.fulfilled;
        const exists = await instance.playerExists(playerId = 1).should.be.fulfilled;
        exists.should.be.equal(true);
    });
});