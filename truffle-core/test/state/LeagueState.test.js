require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeagueState = artifacts.require('LeagueState');

contract('LeagueState', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await LeagueState.new().should.be.fulfilled;
    });

    it('create day state has 0 size', async () => {
        const leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        let count = await instance.leagueStateSize(leagueState).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('append an empty team is valid', async () => {
        const teamState = await instance.teamStateCreate().should.be.fulfilled;
        let leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        leagueState = await instance.leagueStateAppend(leagueState, teamState).should.be.fulfilled;
        let count = await instance.leagueStateSize(leagueState).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        leagueState = await instance.leagueStateAppend(leagueState, teamState).should.be.fulfilled;
        count = await instance.leagueStateSize(leagueState).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('append not empty team state', async () => {
        const playerState = await instance.playerStateCreate(1, 2, 3, 4, 5, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        let leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        leagueState = await instance.leagueStateAppend(leagueState, teamState).should.be.fulfilled;
        let count = await instance.leagueStateSize(leagueState).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        leagueState = await instance.leagueStateAppend(leagueState, teamState).should.be.fulfilled;
        count = await instance.leagueStateSize(leagueState).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('valid league state', async () => {
        const playerState = await instance.playerStateCreate(1, 2, 3, 4, 5, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let result = await instance.isValidLeagueState([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([playerState, 0]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([0]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidLeagueState([playerState,0,0]).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('get team from league state', async () => {
        const leagueState = [2, 3, 0, 4, 2, 1, 0, 4, 5, 0, 2, 0]
        let state = await instance.leagueStateAt(leagueState, 0).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(3);
        state = await instance.leagueStateAt(leagueState, 1).should.be.fulfilled;
        state.length.should.be.equal(3);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(2);
        state[2].toNumber().should.be.equal(1);
        state = await instance.leagueStateAt(leagueState, 2).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(5);
        state = await instance.leagueStateAt(leagueState, 3).should.be.fulfilled;
        state.length.should.be.equal(1);
        state[0].toNumber().should.be.equal(2);
    });

    it('update team state', async () => {
        const playerState = await instance.playerStateCreate(1, 2, 3, 4, 5, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        let leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        leagueState = await instance.leagueStateAppend(leagueState, teamState).should.be.fulfilled;
        leagueState = await instance.leagueStateAppend(leagueState, teamState).should.be.fulfilled;
        let newTeamState = await instance.teamStateCreate().should.be.fulfilled;
        await instance.leagueStateUpdate(leagueState, 1, newTeamState).should.be.rejected;
        const newPlayerState = await instance.playerStateCreate(5, 4, 3, 2, 1, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        newTeamState = await instance.teamStateAppend(newTeamState, newPlayerState).should.be.fulfilled;
        const updatedleagueState = await instance.leagueStateUpdate(leagueState, 1, newTeamState).should.be.fulfilled;
        const resultTeamState = await instance.leagueStateAt(updatedleagueState, 1).should.be.fulfilled;
        const resultPlayerState = await instance.teamStateAt(resultTeamState, 0).should.be.fulfilled;
        resultPlayerState.toString().should.be.equal(newPlayerState.toString());
    });
});