const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
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

    it('is null player virtual', async () => {
        await instance.isVirtual(0).should.be.rejected;
    });

    it('is unexistent player virtual', async () => {
        await instance.isVirtual(1).should.be.rejected;
    });

    it('is existent player virtual', async () => {
        await instance.addTeam("Barca").should.be.fulfilled;
        await instance.isVirtual(1).should.eventually.equal(true);
    });

    it('is existent non virtual player', async () => {
        await instance.setPlayerState(playerId = 1, state = 4).should.be.rejected;
        await instance.addTeam("Barca").should.be.fulfilled;
        await instance.setPlayerState(playerId = 1, state = 4).should.be.fulfilled;
        await instance.isVirtual(playerId = 1).should.eventually.equal(false);
    });

});