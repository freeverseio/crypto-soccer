require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeaguesState = artifacts.require('LeaguesState');

contract('LeaguesState', (accounts) => {
    let instance = null;
    let divider = null;

    beforeEach(async () => {
        instance = await LeaguesState.new().should.be.fulfilled;
        divider = await instance.DIVIDER().should.be.fulfilled;
    });

    it('unexistent league', async () => {
        await leagues.getFinalTeamStateHashes(id).should.be.rejected;
        await leagues.getInitStateHash(id).should.be.rejected;
    });

    it('default hashes values on create league', async () =>{
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const initHash = await leagues.getInitHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await leagues.getFinalTeamStateHashes(id).should.be.fulfilled;
        finalHashes.length.should.be.equal(0);
    });

    it('valid state', async () => {
        let result = await instance.isValid([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([2]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([2, 3, divider, 4, divider, 4]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValid([2, divider, divider, 1]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValid([divider]).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('append an empty team', async () => {
        let result = await instance.append([], []).should.be.fulfilled;
        result.length.should.be.equal(0);
        result = await instance.append([2], []).should.be.fulfilled;
        result.length.should.be.equal(1);
        result[0].toNumber().should.be.equal(2);
    });

    it('append team to league state', async () => {
        let state = await instance.append([], [2]).should.be.fulfilled;
        state.length.should.be.equal(1);
        state[0].toNumber().should.be.equal(2);
        state = await instance.append(state, [3, 4]).should.be.fulfilled;
        state.length.should.be.equal(4);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(0);
        state[2].toNumber().should.be.equal(3);
        state[3].toNumber().should.be.equal(4);
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
        await instance.countTeamPlayers([divider, 2], 0).should.be.rejected;
        const leagueState = [2,3,0,4,2,1,0,4,5,0,2]
        let count = await instance.countTeamPlayers(leagueState, 0).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        count = await instance.countTeamPlayers(leagueState, 1).should.be.fulfilled;
        count.toNumber().should.be.equal(3);
        count = await instance.countTeamPlayers(leagueState, 2).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        count = await instance.countTeamPlayers(leagueState, 3).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('get team from league state', async () => {
        const leagueState = [2, 3, 0, 4, 2, 1, 0, 4, 5, 0, 2]
        let state = await instance.getTeam(leagueState, 0).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(2);
        state[1].toNumber().should.be.equal(3);
        state = await instance.getTeam(leagueState, 1).should.be.fulfilled;
        state.length.should.be.equal(3);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(2);
        state[2].toNumber().should.be.equal(1);
        state = await instance.getTeam(leagueState, 2).should.be.fulfilled;
        state.length.should.be.equal(2);
        state[0].toNumber().should.be.equal(4);
        state[1].toNumber().should.be.equal(5);
        state = await instance.getTeam(leagueState, 3).should.be.fulfilled;
        state.length.should.be.equal(1);
        state[0].toNumber().should.be.equal(2);
    });
});