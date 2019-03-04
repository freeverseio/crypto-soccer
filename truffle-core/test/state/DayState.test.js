require('chai')
    .use(require('chai-as-promised'))
    .should();

const DayState = artifacts.require('DayState');

contract('DayState', (accounts) => {
    let instance = null;
    let TEAMSTATEDIVIDER = null;
    let LEAGUESTATEDIVIDER = null;

    beforeEach(async () => {
        instance = await DayState.new().should.be.fulfilled;
        TEAMSTATEDIVIDER = await instance.TEAMSTATEDIVIDER().should.be.fulfilled;
        LEAGUESTATEDIVIDER = await instance.LEAGUESTATEDIVIDER().should.be.fulfilled;
    });

    it('valid league state', async () => {
        let result = await instance.isValidLeagueState([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([2]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([2, 3, TEAMSTATEDIVIDER, 4, TEAMSTATEDIVIDER, 4]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([2, TEAMSTATEDIVIDER, TEAMSTATEDIVIDER, 1]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidLeagueState([TEAMSTATEDIVIDER]).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('count teams into league state', async () => {
        let count = await instance.countTeamsInState([]).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
        count = await instance.countTeamsInState([2]).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        count = await instance.countTeamsInState([2, 3, 4, 5, 0, 5, 4, 0, 2]).should.be.fulfilled;
        count.toNumber().should.be.equal(3);
    });

    it('count teams into invalid league state', async () => {
        await instance.countTeamsInState([0]).should.be.rejected;
        await instance.countTeamsInState([0,3]).should.be.rejected;
        await instance.countTeamsInState([3,0]).should.be.rejected;
        await instance.countTeamsInState([3,0,0,2]).should.be.rejected;
    });

    it('count players in team', async () => {
        await instance.countTeamPlayers([], 0).should.be.rejected;
        await instance.countTeamPlayers([2], 1).should.be.rejected;
        await instance.countTeamPlayers([TEAMSTATEDIVIDER, 2], 0).should.be.rejected;
        const dayState = [2,3,0,4,2,1,0,4,5,0,2]
        let count = await instance.countTeamPlayers(dayState, 0).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        count = await instance.countTeamPlayers(dayState, 1).should.be.fulfilled;
        count.toNumber().should.be.equal(3);
        count = await instance.countTeamPlayers(dayState, 2).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        count = await instance.countTeamPlayers(dayState, 3).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('get team from league state', async () => {
        const dayState = [2, 3, 0, 4, 2, 1, 0, 4, 5, 0, 2]
        let state = await instance.getTeam(dayState, 0).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(3);
        state = await instance.getTeam(dayState, 1).should.be.fulfilled;
        state.length.should.be.equal(3);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(2);
        state[2].toNumber().should.be.equal(1);
        state = await instance.getTeam(dayState, 2).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(5);
        state = await instance.getTeam(dayState, 3).should.be.fulfilled;
        state.length.should.be.equal(1);
        state[0].toNumber().should.be.equal(2);
    });

    it('append team state to league state', async () => {
        let dayState = await instance.dayStateCreate().should.be.fulfilled;
        dayState = await instance.dayStateAppend(dayState, [4,5,6,7]).should.be.fulfilled;
        dayState.length.should.be.equal(4);
        await instance.dayStateAppend(dayState, []).should.be.rejected;
        dayState = await instance.dayStateAppend(dayState, [2]).should.be.fulfilled;
        dayState.length.should.be.equal(6);
        dayState[0].toNumber().should.be.equal(4);
        dayState[1].toNumber().should.be.equal(5);
        dayState[2].toNumber().should.be.equal(6);
        dayState[3].toNumber().should.be.equal(7);
        dayState[4].toNumber().should.be.equal(0);
        dayState[5].toNumber().should.be.equal(2);
    });
});