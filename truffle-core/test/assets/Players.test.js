const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Players = artifacts.require('PlayersMock');
const PlayerStateLib = artifacts.require('PlayerState');

contract('Players', (accounts) => {
    let players = null;
    let playerStateLib = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        players = await Players.new(playerStateLib.address).should.be.fulfilled;
    });

    it('query null player id', async () => {
        await players.getPlayerTeam(0).should.be.rejected;
    });

    it('query non created player id', async () => {
        await players.getPlayerTeam(1).should.be.rejected;
    });

    it('get player team of existing player', async () => {
        const nPLayersPerTeam = await players.PLAYERS_PER_TEAM().should.be.fulfilled;
        await players.addTeam("Barca").should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const teamId = await players.getPlayerTeam(playerId).should.be.fulfilled;
            teamId.toNumber().should.be.equal(1);
        }
        await players.getPlayerTeam(nPLayersPerTeam+1).should.be.rejected;
    });

    it('computed skills with rnd = 0 is 50 each', async () => {
        let skills = await players.computeSkills(0).should.be.fulfilled;
        skills.forEach(skill => (skill.toNumber().should.be.equal(50)));
    });

    it('int hash is deterministic', async () => {
        const rand0 = await players.intHash("Barca0").should.be.fulfilled;
        const rand1 = await players.intHash("Barca0").should.be.fulfilled;
        rand0.should.be.bignumber.equal(rand1);
        const rand2 = await players.intHash("Barca1").should.be.fulfilled;
        rand0.should.be.bignumber.not.equal(rand2);
        rand0.should.be.bignumber.equal('64856073772839990506814373782217928521534618466099710722049665631602958392435');
    });

    it('sum of computed skills is 250', async () => {
        for (let i = 0; i < 10; i++) {
            const seed = await players.intHash("Barca" + i).should.be.fulfilled;
            const skills = await players.computeSkills(seed).should.be.fulfilled;
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            sum.should.be.equal(250);
        }
    });

    it('get player pos in team', async () => {
        const nPLayersPerTeam = await players.PLAYERS_PER_TEAM().should.be.fulfilled;
        await players.addTeam("Barca").should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const pos = await players.getPlayerPosInTeam(playerId).should.be.fulfilled;
            pos.toNumber().should.be.equal(playerId - 1);
        }
        await players.getPlayerPosInTeam(nPLayersPerTeam+1).should.be.rejected;
    })

    it('get existing player skills', async () => {
        const numSkills = await players.NUM_SKILLS().should.be.fulfilled;
        await players.addTeam("Barca").should.be.fulfilled;
        const skills = await players.getPlayerSkills(playerId = 10).should.be.fulfilled;
        skills.length.should.be.equal(numSkills.toNumber());
        skills[0].should.be.bignumber.equal('48');
        skills[1].should.be.bignumber.equal('72');
        skills[2].should.be.bignumber.equal('51');
        skills[3].should.be.bignumber.equal('42');
        skills[4].should.be.bignumber.equal('37');
        const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
        sum.should.be.equal(250);
    });

    it('compute player birth', async () => {
        await players.addTeam("Barca").should.be.fulfilled;
        const birth = await players.computeBirth(0).should.be.fulfilled;
        birth.should.be.bignumber.equal('406');
    });

    // it('exchange players team', async () => {
    //     await players.addTeam("Barca").should.be.fulfilled;
    //     await players.addTeam("Madrid").should.be.fulfilled;
    //     await players.exchangePlayersTeams(playerId0 = 8, playerId1 = 17).should.be.fulfilled;
    //     const teamPlayer0 = await players.getPlayerTeam(playerId0).should.be.fulfilled;
    //     teamPlayer0.should.be.bignumber.equal('2');
    //     const teamPlayer1 = await players.getPlayerTeam(playerId1).should.be.fulfilled;
    //     teamPlayer1.should.be.bignumber.equal('1');
    // });

    it('get non virtual player team', async () => {
        await players.addTeam("Barca").should.be.fulfilled;
        const teamBefore = await players.getPlayerTeam(playerId = 1).should.be.fulfilled;
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
        await players.setPlayerState(state).should.be.fulfilled;
        const teamAfter = await players.getPlayerTeam(playerId = 1).should.be.fulfilled;
        teamAfter.should.be.bignumber.not.equal(teamBefore);
        teamAfter.should.be.bignumber.equal('4');
    });
});
 