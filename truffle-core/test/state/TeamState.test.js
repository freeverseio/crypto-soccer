require('chai')
    .use(require('chai-as-promised'))
    .should();

const TeamState = artifacts.require('TeamState');

contract('TeamState', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await TeamState.new().should.be.fulfilled;
    });

    it('create team state', async () => {
        const teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState.length.should.be.equal(0);
    });

    it('count players in team state', async () => {
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        let count = await instance.teamStateSize(teamState).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
        const playerState = await instance.playerStateCreate(1, 2, 3, 4, 5, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        count = await instance.teamStateSize(teamState).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    })

    it('append player state to team state', async () => {
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
        const playerState = await instance.playerStateCreate(0,0,0,0,0,0, playerId = 1, 0, 0,0,0,0,0).should.be.fulfilled;
        let result = await instance.isValidTeamState([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidTeamState([playerState]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.isValidTeamState([0]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidTeamState([0, playerState]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidTeamState([playerState, 0]).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('team rating', async () => {
        const nPlayers = 100;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        for (var i = 1; i < nPlayers; i += 5) {
            const playerState = await instance.playerStateCreate(i, i + 1, i + 2, i + 3, i + 4, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
            teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        }
        const rating = await instance.computeTeamRating(teamState).should.be.fulfilled;
        rating.toNumber().should.be.equal(nPlayers * (nPlayers + 1) / 2);
    });

    it('evolve team of delta 0 (1)', async () => {
        const playerState = await instance.playerStateCreate(0,0,0,0,0,0, playerId = 1, 0, 0,0,0,0,0).should.be.fulfilled;
        const teamState = [playerState];
        const evolvedTeamState = await instance.teamStateEvolve(teamState, 0).should.be.fulfilled;
        const valid = await instance.isValidTeamState(evolvedTeamState).should.be.fulfilled;
        valid.should.be.equal(true);
    })

    it('evolve team of delta 0 (2)', async () => {
        const playerState = await instance.playerStateCreate(1, 2, 3, 4, 5, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState).should.be.fulfilled;
        const evolvedTeamState = await instance.teamStateEvolve(teamState, 0).should.be.fulfilled;
        const valid = await instance.isValidTeamState(evolvedTeamState).should.be.fulfilled;
        valid.should.be.equal(true);
        const evolvedPlayerState = await instance.teamStateAt(evolvedTeamState, 0).should.be.fulfilled;
        const defence = await instance.getDefence(evolvedPlayerState).should.be.fulfilled;
        defence.toString().should.be.equal('1');
    });

    it('team state evolve', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        let playerState0 = await instance.playerStateCreate(defence, speed, pass, shoot, endurance, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let playerState1 = await instance.playerStateCreate(defence + 1, speed + 1, pass + 1, shoot + 1, endurance + 1, 0, playerId = 2, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let teamState = await instance.teamStateCreate().should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState0).should.be.fulfilled;
        teamState = await instance.teamStateAppend(teamState, playerState1).should.be.fulfilled;
        const delta = 3;
        teamState = await instance.teamStateEvolve(teamState, delta).should.be.fulfilled;
        playerState0 = await instance.teamStateAt(teamState, 0).should.be.fulfilled;
        let skill = await instance.getDefence(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence); // + delta);
        skill = await instance.getSpeed(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed); //  + delta);
        skill = await instance.getPass(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass); //  + delta);
        skill = await instance.getShoot(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot); // + delta);
        skill = await instance.getEndurance(playerState0).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance); // + delta);
        playerState1 = await instance.teamStateAt(teamState, 1).should.be.fulfilled;
        skill = await instance.getDefence(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence + 1); // + delta);
        skill = await instance.getSpeed(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed + 1); // + delta);
        skill = await instance.getPass(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass + 1); // + delta);
        skill = await instance.getShoot(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot + 1); // + delta);
        skill = await instance.getEndurance(playerState1).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance + 1); // + delta);
    });
});