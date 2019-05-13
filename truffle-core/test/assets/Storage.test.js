const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Storage = artifacts.require('StorageMock');
const PlayerStateLib = artifacts.require('PlayerState');

contract('Storage', (accounts) => {
    let instance = null;
    let playerStateLib = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        instance = await Storage.new(playerStateLib.address).should.be.fulfilled;
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
        await instance.setPlayerState(4).should.be.rejected;
        await instance.addTeam("Barca").should.be.fulfilled;
        const state = await playerStateLib.playerStateCreate(
            defence = 3,
            speed = 3,
            pass = 3,
            shoot = 3,
            endurance = 3,
            monthOfBirthInUnixTime = 3,
            playerId = 1,
            currentTeamId = 4,
            currentShirtNum = 3,
            prevLeagueId = 3,
            prevTeamPosInLeague = 3,
            prevShirtNumInLeague = 3,
            lastSaleBlock = 3
        ).should.be.fulfilled;
        await instance.setPlayerState(state).should.be.fulfilled;
        await instance.isVirtual(playerId = 1).should.eventually.equal(false);
    });

});