require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeagueState = artifacts.require('LeagueState');

contract('LeagueState', (accounts) => {
    let instance = null;
    let divider = null;

    beforeEach(async () => {
        instance = await LeagueState.new().should.be.fulfilled;
        divider = await instance.DIVIDER().should.be.fulfilled;
    });

    it('valid state', async () => {
        let result = await instance.isValid([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([2]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([2, 3, divider, 4, divider, 4]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([2, divider, divider, 1]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([divider]).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('append an empty team', async () => {
        let result = await instance.appendTeamToLeagueState([], []).should.be.fulfilled;
        result.length.should.be.equal(0);
        result = await instance.appendTeamToLeagueState([2], []).should.be.fulfilled;
        result.length.should.be.equal(1);
        result[0].toNumber().should.be.equal(2);
    });

    it('append team to league state', async () => {
        let state = await instance.appendTeamToLeagueState([], [2]).should.be.fulfilled;
        state.length.should.be.equal(1);
        state[0].toNumber().should.be.equal(2);
        state = await instance.appendTeamToLeagueState(state, [3, 4]).should.be.fulfilled;
        state.length.should.be.equal(4);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(0);
        state[2].toNumber().should.be.equal(3);
        state[3].toNumber().should.be.equal(4);
    });

    it('count team states into league state', async () => {
        let count = await instance.countTeams([]).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
        count = await instance.countTeams([2]).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        count = await instance.countTeams([2, 3, 4, 5, 0, 5, 4, 0, 2]).should.be.fulfilled;
        count.toNumber().should.be.equal(3);
    });

    it('count team states into invalid league state', async () => {
        await instance.countTeams([0]).should.be.rejected;
        await instance.countTeams([0,3]).should.be.rejected;
        await instance.countTeams([3,0]).should.be.rejected;
    });
});