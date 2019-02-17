require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeagueState = artifacts.require('LeagueState');

contract('LeagueState', (accounts) => {
    let leagueState = null;

    beforeEach(async () => {
        leagueState = await LeagueState.new().should.be.fulfilled;
    });

    it('append team to league state', async () => {
        await leagueState.appendTeamToLeagueState([], []).should.be.rejected;
        let state = await leagueState.appendTeamToLeagueState([], [2]).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(0);
        state = await leagueState.appendTeamToLeagueState(state, [3, 4]).should.be.fulfilled;
        state.length.should.be.equal(5);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(0);
        state[2].toNumber().should.be.equal(3);
        state[3].toNumber().should.be.equal(4);
        state[4].toNumber().should.be.equal(0);
    });
});