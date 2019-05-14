const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Storage = artifacts.require('StorageMock');
const PlayerStateLib = artifacts.require('PlayerState');
 
contract('Storage', (accounts) => {
    let instance = null;
    let playerStateLib = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        instance = await Storage.new(playerStateLib.address).should.be.fulfilled;
    });

    it('get playerIds of the team', async () => {
        await instance.addTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        let playerIds = await instance.getTeamPlayerIds(1).should.be.fulfilled;
        playerIds.length.should.be.equal(11);
        for (let pos = 0; pos < 11 ; pos++) 
            playerIds[pos].should.be.bignumber.equal((pos+1).toString());

        await instance.addTeam(name = "Madrid",accounts[1]).should.be.fulfilled;
        await instance.exchangePlayersTeams(playerId0 = 11, playerId1 = 14).should.be.fulfilled;
        playerIds = await instance.getTeamPlayerIds(1).should.be.fulfilled;
        playerIds.length.should.be.equal(11);
        for (let pos = 0; pos < 10 ; pos++) 
            playerIds[pos].should.be.bignumber.equal((pos+1).toString());
        playerIds[10].should.be.bignumber.equal('14');
    });

    it('add team with different owner than the sender', async () => {
        await instance.addTeam('Barca', accounts[1]).should.be.fulfilled;
        const owner = await instance.getTeamOwner('Barca').should.be.fulfilled;
        owner.should.be.equal(accounts[1]);
    })

    it('add 2 teams with same name', async() => {
        await instance.addTeam('Barca', accounts[1]).should.be.fulfilled;
        await instance.addTeam('Barca', accounts[1]).should.be.rejected;
    })

    it('team exists', async () => {
        let result = await instance.teamExists(0).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.teamExists(1).should.be.fulfilled;
        result.should.be.equal(false);
        await instance.addTeam("Barca", accounts[1]).should.be.fulfilled;
        result = await instance.teamExists(1).should.be.fulfilled;
        result.should.be.equal(true);
        result = await instance.teamExists(2).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('initial number of team', async () => {
        const count = await instance.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('get name of invalid team', async () => {
        await instance.getTeamName(0).should.be.rejected;
    });

    it('get name of unexistent team', async () => {
        await instance.getTeamName(1).should.be.rejected;
    });

    it('existence of null player', async () => {
        const exists = await instance.playerExists(playerId = 0).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('existence of unexistent player', async () => {
        const exists = await instance.playerExists(playerId = 1).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('existence of existent player', async () => {
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        const exists = await instance.playerExists(playerId = 1).should.be.fulfilled;
        exists.should.be.equal(true);
    });

    it('is null player virtual', async () => {
        await instance.isVirtual(0).should.be.rejected;
    });

    it('is unexistent player virtual', async () => {
        await instance.isVirtual(1).should.be.rejected;
    });

    it('is existent player virtual', async () => {
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        await instance.isVirtual(1).should.eventually.equal(true);
    });

    it('is existent non virtual player', async () => {
        await instance.setPlayerState(4).should.be.rejected;
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        const state = await playerStateLib.playerStateCreate(
            defence = 3,
            speed = 3,
            pass = 3,
            shoot = 3,
            endurance = 3,
            monthOfBirthInUnixTime = 3,
            playerId = 1,
            currentTeamId = 1,
            currentShirtNum = 3,
            prevLeagueId = 3,
            prevTeamPosInLeague = 3,
            prevShirtNumInLeague = 3,
            lastSaleBlock = 3
        ).should.be.fulfilled;
        await instance.setPlayerState(state).should.be.fulfilled;
        await instance.isVirtual(playerId = 1).should.eventually.equal(false);
    });

    it('get state of virtual player', async () => {
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        const state = await instance.getPlayerState(playerId = 1).should.be.fulfilled;
        state.should.be.bignumber.equal('473533131866555579417557877411906949081105664487195487081826231992180539392');
    });

    it('exchange players team', async () => {
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        await instance.addTeam("Madrid",accounts[1]).should.be.fulfilled;
        await instance.exchangePlayersTeams(playerId0 = 8, playerId1 = 19).should.be.fulfilled;
        const teamPlayer0 = await instance.getPlayerTeam(playerId0).should.be.fulfilled;
        teamPlayer0.should.be.bignumber.equal('2');
        const teamPlayer1 = await instance.getPlayerTeam(playerId1).should.be.fulfilled;
        teamPlayer1.should.be.bignumber.equal('1');
    });

    it('query null player id', async () => {
        await instance.getPlayerTeam(0).should.be.rejected;
    });

    it('query non created player id', async () => {
        await instance.getPlayerTeam(1).should.be.rejected;
    });

    it('get player team of existing player', async () => {
        const nPLayersPerTeam = await instance.PLAYERS_PER_TEAM().should.be.fulfilled;
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const teamId = await instance.getPlayerTeam(playerId).should.be.fulfilled;
            teamId.toNumber().should.be.equal(1);
        }
        await instance.getPlayerTeam(nPLayersPerTeam+1).should.be.rejected;
    });

    it('computed skills with rnd = 0 is 50 each', async () => {
        let skills = await instance.computeSkills(0).should.be.fulfilled;
        skills.forEach(skill => (skill.toNumber().should.be.equal(50)));
    });

    it('int hash is deterministic', async () => {
        const rand0 = await instance.intHash("Barca0").should.be.fulfilled;
        const rand1 = await instance.intHash("Barca0").should.be.fulfilled;
        rand0.should.be.bignumber.equal(rand1);
        const rand2 = await instance.intHash("Barca1").should.be.fulfilled;
        rand0.should.be.bignumber.not.equal(rand2);
        rand0.should.be.bignumber.equal('64856073772839990506814373782217928521534618466099710722049665631602958392435');
    });

    it('sum of computed skills is 250', async () => {
        for (let i = 0; i < 10; i++) {
            const seed = await instance.intHash("Barca" + i).should.be.fulfilled;
            const skills = await instance.computeSkills(seed).should.be.fulfilled;
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            sum.should.be.equal(250);
        }
    });

    it('get player pos in team', async () => {
        const nPLayersPerTeam = await instance.PLAYERS_PER_TEAM().should.be.fulfilled;
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const pos = await instance.getPlayerPosInTeam(playerId).should.be.fulfilled;
            pos.toNumber().should.be.equal(playerId - 1);
        }
        await instance.getPlayerPosInTeam(nPLayersPerTeam+1).should.be.rejected;
    })

    it('get existing virtual player skills', async () => {
        const numSkills = await instance.NUM_SKILLS().should.be.fulfilled;
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        const skills = await instance.getPlayerSkills(playerId = 10).should.be.fulfilled;
        skills.length.should.be.equal(numSkills.toNumber());
        skills[0].should.be.bignumber.equal('48');
        skills[1].should.be.bignumber.equal('72');
        skills[2].should.be.bignumber.equal('51');
        skills[3].should.be.bignumber.equal('42');
        skills[4].should.be.bignumber.equal('37');
        const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
        sum.should.be.equal(250);
    });

    it('get existing non virtual player skills', async () => {
        const numSkills = await instance.NUM_SKILLS().should.be.fulfilled;
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        const state = await playerStateLib.playerStateCreate(
            defence = 1,
            speed = 2,
            pass = 3,
            shoot = 4,
            endurance = 5,
            monthOfBirthInUnixTime = 6,
            playerId = 10,
            currentTeamId = 1,
            currentShirtNum = 3,
            prevLeagueId = 3,
            prevTeamPosInLeague = 3,
            prevShirtNumInLeague = 3,
            lastSaleBlock = 3
        ).should.be.fulfilled;
        await instance.setPlayerState(state).should.be.fulfilled;
        const skills = await instance.getPlayerSkills(playerId = 10).should.be.fulfilled;
        skills[0].should.be.bignumber.equal('1');
        skills[1].should.be.bignumber.equal('2');
        skills[2].should.be.bignumber.equal('3');
        skills[3].should.be.bignumber.equal('4');
        skills[4].should.be.bignumber.equal('5');
    });

    it('compute player birth', async () => {
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        const birth = await instance.computeBirth(0).should.be.fulfilled;
        birth.should.be.bignumber.equal('406');
    });

    it('get non virtual player team', async () => {
        await instance.addTeam("Barca",accounts[1]).should.be.fulfilled;
        await instance.addTeam("Madrid",accounts[1]).should.be.fulfilled;
        const teamBefore = await instance.getPlayerTeam(playerId = 1).should.be.fulfilled;
        const state = await playerStateLib.playerStateCreate(
            defence = 3,
            speed = 3,
            pass = 3,
            shoot = 3,
            endurance = 3,
            monthOfBirthInUnixTime = 3,
            playerId = 1,
            currentTeamId = 2,
            currentShirtNum = 3,
            prevLeagueId = 3,
            prevTeamPosInLeague = 3,
            prevShirtNumInLeague = 3,
            lastSaleBlock = 3
        ).should.be.fulfilled;
        await instance.setPlayerState(state).should.be.fulfilled;
        const teamAfter = await instance.getPlayerTeam(playerId = 1).should.be.fulfilled;
        teamAfter.should.be.bignumber.not.equal(teamBefore);
        teamAfter.should.be.bignumber.equal('2');
    });

    it('create team', async () => {
        const receipt = await instance.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        const count = await instance.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        let teamName = receipt.logs[0].args.teamName;
        teamName.should.be.equal("Barca",accounts[1]);
        const teamId = receipt.logs[0].args.teamId.toNumber();
        teamId.should.be.equal(1);
        teamName = await instance.getTeamName(teamId).should.be.fulfilled;
        teamName.should.be.equal("Barca",accounts[1]);
    });

    it('get playersId from teamId and pos in team', async () => {
        await instance.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.rejected;
        await instance.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        await instance.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=11).should.be.rejected;
        let playerId = await instance.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.fulfilled;
        playerId.toNumber().should.be.equal(1);
        playerId = await instance.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=10).should.be.fulfilled;
        playerId.toNumber().should.be.equal(11);
    });

    it('sign team to league', async () => {
        await instance.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
        await instance.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        await instance.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.fulfilled;
        const currentHistory = await instance.getTeamCurrentHistory(1).should.be.fulfilled;
        currentHistory.currentLeagueId.should.be.bignumber.equal('1');
        currentHistory.posInCurrentLeague.should.be.bignumber.equal('3');
        currentHistory.prevLeagueId.should.be.bignumber.equal('0');
        currentHistory.posInPrevLeague.should.be.bignumber.equal('0');
    });
});