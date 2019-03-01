require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeaguesState = artifacts.require('LeaguesState');

contract('LeaguesState', (accounts) => {
    let instance = null;
    let divider = null;
    const initBlock = 1;
    const step = 1;
    const id = 0;
    const teamIds = [1, 2];

    beforeEach(async () => {
        instance = await LeaguesState.new().should.be.fulfilled;
        divider = await instance.DIVIDER().should.be.fulfilled;
    });

    it('unexistent league', async () => {
        await instance.getFinalTeamStateHashes(id).should.be.rejected;
        await instance.getInitStateHash(id).should.be.rejected;
    });

    it('default hashes values on create league', async () =>{
        await instance.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const initHash = await instance.getInitStateHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await instance.getFinalTeamStateHashes(id).should.be.fulfilled;
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

    it('create player state', async () => {
        const defence = 3;
        const speed = 23;
        const pass = 2;
        const shoot = 21;
        const endurance = 10;
        const state = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
        (state.toNumber() & 0xff).should.be.equal(endurance);
        (state.toNumber() >> 8 & 0xff).should.be.equal(shoot);
        (state.toNumber() >> 8*2 & 0xff).should.be.equal(pass);
        (state.toNumber() >> 8*3 & 0xff).should.be.equal(speed);
        (state.toNumber() >> 8*4 & 0xff).should.be.equal(endurance);
    });

    it('create team state', async () => {
        const teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState.length.should.be.equal(0);
    });

    it('leagueStateAppend player state to team state', async () => {
        const playerState0 = 0x546ab;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState0).should.be.fulfilled;
        teamState.length.should.be.equal(1);
        const playerState1 = 0x435;
        teamState = await instance.teamStateAppend(teamState, playerState1).should.be.fulfilled;
        teamState.length.should.be.equal(2);
        teamState[0].toNumber().should.be.equal(playerState0);
        teamState[1].toNumber().should.be.equal(playerState1);
    });

    it('valid team state', async () => {
        let result = await instance.isValidTeamState([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidTeamState([0]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidTeamState([0, 4,3,2,1]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidTeamState([8,0,34]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidTeamState([4,0]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidTeamState([3,4,5,76]).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('leagueStateAppend team state to league state', async () => {
        let leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        leagueState = await instance.leagueStateAppend(leagueState, [4,5,6,7]).should.be.fulfilled;
        leagueState.length.should.be.equal(4);
        await instance.leagueStateAppend(leagueState, []).should.be.rejected;
        leagueState = await instance.leagueStateAppend(leagueState, [2]).should.be.fulfilled;
        leagueState.length.should.be.equal(6);
        leagueState[0].toNumber().should.be.equal(4);
        leagueState[1].toNumber().should.be.equal(5);
        leagueState[2].toNumber().should.be.equal(6);
        leagueState[3].toNumber().should.be.equal(7);
        leagueState[4].toNumber().should.be.equal(0);
        leagueState[5].toNumber().should.be.equal(2);
    });

    it('team rating', async () => {
        const nPlayers = 100;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        for (var i = 1; i < nPlayers; i += 5) {
            const playerState = await instance.playerStateCreate(i, i + 1, i + 2, i + 3, i + 4).should.be.fulfilled;
            teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        }
        const rating = await instance.computeTeamRating(teamState).should.be.fulfilled;
        rating.toNumber().should.be.equal(nPlayers * (nPlayers + 1) / 2);
    });

    it('is valid player state', async () => {
        let result = await instance.isValidPlayerState(-1).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidPlayerState(0).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('skills getters from state player', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        const playerState = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
        let result = await instance.getDefence(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(defence);
        result = await instance.getSpeed(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(speed);
        result = await instance.getPass(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(pass);
        result = await instance.getShoot(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(shoot);
        result = await instance.getEndurance(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(endurance);
    });

    it('player state evolve', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        const playerState = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
        const delta = 3;
        const updatedState = await instance.playerStateEvolve(playerState, delta).should.be.fulfilled;
        updatedState.toNumber().should.not.be.equal(playerState.toNumber());
        let skill = await instance.getDefence(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence + delta);
        skill = await instance.getSpeed(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed + delta);
        skill = await instance.getPass(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass + delta);
        skill = await instance.getShoot(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot + delta);
        skill = await instance.getEndurance(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance + delta);
    });

    it('team state evolve', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        let playerState0 = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
        let playerState1 = await instance.playerStateCreate(defence+1, speed+1, pass+1, shoot+1, endurance+1).should.be.fulfilled;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState0).should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState1).should.be.fulfilled;
        const delta = 3;
        teamState = await instance.teamStateEvolve(teamState, delta).should.be.fulfilled;
        playerState0 = await instance.getPlayerState(teamState, 0).should.be.fulfilled;
        let skill = await instance.getDefence(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence + delta);
        skill = await instance.getSpeed(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed + delta);
        skill = await instance.getPass(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass + delta);
        skill = await instance.getShoot(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot + delta);
        skill = await instance.getEndurance(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance + delta);
        playerState1 = await instance.getPlayerState(teamState, 1).should.be.fulfilled;
        skill = await instance.getDefence(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence + 1 + delta);
        skill = await instance.getSpeed(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed + 1 + delta);
        skill = await instance.getPass(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass + 1 + delta);
        skill = await instance.getShoot(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot + 1 + delta);
        skill = await instance.getEndurance(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance + 1 + delta);
    });
});