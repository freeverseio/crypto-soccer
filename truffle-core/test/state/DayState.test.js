require('chai')
    .use(require('chai-as-promised'))
    .should();

const DayState = artifacts.require('DayState');

contract('DayState', (accounts) => {
    let instance = null;
    let TEAMSTATEDIVIDER = null;

    beforeEach(async () => {
        instance = await DayState.new().should.be.fulfilled;
        TEAMSTATEDIVIDER = await instance.TEAMSTATEDIVIDER().should.be.fulfilled;
    });

    it('create day state has 0 size', async () => {
        const dayState = await instance.dayStateCreate().should.be.fulfilled;
        let count = await instance.dayStateSize(dayState).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('append an empty team is valid', async () => {
        const teamState = await instance.teamStateCreate().should.be.fulfilled;
        let dayState = await instance.dayStateCreate().should.be.fulfilled;
        dayState = await instance.dayStateAppend(dayState, teamState).should.be.fulfilled;
        let count = await instance.dayStateSize(dayState).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        dayState = await instance.dayStateAppend(dayState, teamState).should.be.fulfilled;
        count = await instance.dayStateSize(dayState).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('append not empty team state', async () => {
        const playerState = await instance.playerStateCreate(1, 2, 3, 4, 5).should.be.fulfilled;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        let dayState = await instance.dayStateCreate().should.be.fulfilled;
        dayState = await instance.dayStateAppend(dayState, teamState).should.be.fulfilled;
        let count = await instance.dayStateSize(dayState).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        dayState = await instance.dayStateAppend(dayState, teamState).should.be.fulfilled;
        count = await instance.dayStateSize(dayState).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('valid league state', async () => {
        let result = await instance.isValidLeagueState([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([2]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidLeagueState([2, 3, TEAMSTATEDIVIDER, 4, TEAMSTATEDIVIDER, 4]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidLeagueState([2, TEAMSTATEDIVIDER, TEAMSTATEDIVIDER, 1, TEAMSTATEDIVIDER]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([TEAMSTATEDIVIDER]).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('count players in team', async () => {
        await instance.countTeamPlayers([], 0).should.be.rejected;
        await instance.countTeamPlayers([2], 1).should.be.rejected;
        await instance.countTeamPlayers([TEAMSTATEDIVIDER, 2], 0).should.be.rejected;
        const dayState = [2,3,0,4,2,1,0,4,5,0,2,0]
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
        const dayState = [2, 3, 0, 4, 2, 1, 0, 4, 5, 0, 2, 0]
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
});