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
        let count = await instance.countLeagueStateTeams([]).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
        count = await instance.countLeagueStateTeams([2]).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        count = await instance.countLeagueStateTeams([2, 3, 4, 5, 0, 5, 4, 0, 2]).should.be.fulfilled;
        count.toNumber().should.be.equal(3);
    });

    it('count team states into invalid league state', async () => {
        await instance.countLeagueStateTeams([0]).should.be.rejected;
        await instance.countLeagueStateTeams([0,3]).should.be.rejected;
        await instance.countLeagueStateTeams([3,0]).should.be.rejected;
    });

    it('count players in team', async () => {
        await instance.countPlayersInTeam([], 0).should.be.rejected;
        await instance.countPlayersInTeam([2], 1).should.be.rejected;
        await instance.countPlayersInTeam([divider, 2], 0).should.be.rejected;
        const leagueState = [2,3,0,4,2,1,0,4,5,0,2]
        let count = await instance.countPlayersInTeam(leagueState, 0).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        count = await instance.countPlayersInTeam(leagueState, 1).should.be.fulfilled;
        count.toNumber().should.be.equal(3);
        count = await instance.countPlayersInTeam(leagueState, 2).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        count = await instance.countPlayersInTeam(leagueState, 3).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('get team state from league state', async () => {
        const leagueState = [2, 3, 0, 4, 2, 1, 0, 4, 5, 0, 2]
        let state = await instance.getTeamState(leagueState, 0).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(3);
        state = await instance.getTeamState(leagueState, 1).should.be.fulfilled;
        state.length.should.be.equal(3);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(2);
        state[2].toNumber().should.be.equal(1);
        state = await instance.getTeamState(leagueState, 2).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(5);
        state = await instance.getTeamState(leagueState, 3).should.be.fulfilled;
        state.length.should.be.equal(1);
        state[0].toNumber().should.be.equal(2);
    });
});