const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Assets = artifacts.require('AssetsMock');
const PlayerStateLib = artifacts.require('PlayerState');
 
contract('Assets', (accounts) => {
    let assets = null;
    let playerStateLib = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
    });

    it('get playerIds of the team', async () => {
        await assets.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        let playerIds = await assets.getTeamPlayerIds(1).should.be.fulfilled;
        playerIds.length.should.be.equal(11);
        for (let pos = 0; pos < 11 ; pos++) 
            playerIds[pos].should.be.bignumber.equal((pos+1).toString());

        await assets.createTeam(name = "Madrid",accounts[1]).should.be.fulfilled;
        await assets.exchangePlayersTeams(playerId0 = 11, playerId1 = 14).should.be.fulfilled;
        playerIds = await assets.getTeamPlayerIds(1).should.be.fulfilled;
        playerIds.length.should.be.equal(11);
        for (let pos = 0; pos < 10 ; pos++) 
            playerIds[pos].should.be.bignumber.equal((pos+1).toString());
        playerIds[10].should.be.bignumber.equal('14');
    });

    it('add team with different owner than the sender', async () => {
        await assets.createTeam('Barca', accounts[1]).should.be.fulfilled;
        const owner = await assets.getTeamOwner('Barca').should.be.fulfilled;
        owner.should.be.equal(accounts[1]);
    })

    it('add 2 teams with same name', async() => {
        await assets.createTeam('Barca', accounts[1]).should.be.fulfilled;
        await assets.createTeam('Barca', accounts[1]).should.be.rejected;
    })

    it('team exists', async () => {
        let result = await assets.teamExists(0).should.be.fulfilled;
        result.should.be.equal(false);
        result = await assets.teamExists(1).should.be.fulfilled;
        result.should.be.equal(false);
        await assets.createTeam("Barca", accounts[1]).should.be.fulfilled;
        result = await assets.teamExists(1).should.be.fulfilled;
        result.should.be.equal(true);
        result = await assets.teamExists(2).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('initial number of team', async () => {
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('get name of invalid team', async () => {
        await assets.getTeamName(0).should.be.rejected;
    });

    it('get name of unexistent team', async () => {
        await assets.getTeamName(1).should.be.rejected;
    });

    it('existence of null player', async () => {
        const exists = await assets.playerExists(playerId = 0).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('existence of unexistent player', async () => {
        const exists = await assets.playerExists(playerId = 1).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('existence of existent player', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        const exists = await assets.playerExists(playerId = 1).should.be.fulfilled;
        exists.should.be.equal(true);
    });

    it('is null player virtual', async () => {
        await assets.isVirtual(0).should.be.rejected;
    });

    it('is unexistent player virtual', async () => {
        await assets.isVirtual(1).should.be.rejected;
    });

    it('is existent player virtual', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        await assets.isVirtual(1).should.eventually.equal(true);
    });

    it('set player state of existent virtual player', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        let state = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
        const currentBlock = 5; // TODO: get it properly
        state = await playerStateLib.setLastSaleBlock(state, currentBlock).should.be.fulfilled;
        await assets.setPlayerState(state).should.be.fulfilled;
        const resultState = await assets.getPlayerState(playerId).should.be.fulfilled;
        resultState.should.be.bignumber.equal(state);
    });

    it('is existent non virtual player', async () => {
        await assets.setPlayerState(4).should.be.rejected;
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
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
        await assets.setPlayerState(state).should.be.fulfilled;
        await assets.isVirtual(playerId = 1).should.eventually.equal(false);
    });

    it('get state of virtual player', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        const state = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
        state.should.be.bignumber.equal('473533131866555579417557877411906949081105664487195487081826231992180539392');
    });

    it('exchange players team', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        await assets.createTeam("Madrid",accounts[1]).should.be.fulfilled;
        await assets.exchangePlayersTeams(playerId0 = 8, playerId1 = 19).should.be.fulfilled;
        const statePlayer0 = await assets.getPlayerState(playerId0).should.be.fulfilled;
        const teamPlayer0 = await playerStateLib.getCurrentTeamId(statePlayer0).should.be.fulfilled;
        teamPlayer0.should.be.bignumber.equal('2');
        const statePlayer1 = await assets.getPlayerState(playerId1).should.be.fulfilled;
        const teamPlayer1 = await playerStateLib.getCurrentTeamId(statePlayer1).should.be.fulfilled;
        teamPlayer1.should.be.bignumber.equal('1');
    });

    it('get player state of existing player', async () => {
        const nPLayersPerTeam = await assets.PLAYERS_PER_TEAM().should.be.fulfilled;
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++)
            await assets.getPlayerState(playerId).should.be.fulfilled;
        await assets.getPlayerState(nPLayersPerTeam+1).should.be.rejected;
    });

    it('computed skills with rnd = 0 is 50 each', async () => {
        let skills = await assets.computeSkills(0).should.be.fulfilled;
        skills.forEach(skill => (skill.toNumber().should.be.equal(50)));
    });

    it('int hash is deterministic', async () => {
        const rand0 = await assets.intHash("Barca0").should.be.fulfilled;
        const rand1 = await assets.intHash("Barca0").should.be.fulfilled;
        rand0.should.be.bignumber.equal(rand1);
        const rand2 = await assets.intHash("Barca1").should.be.fulfilled;
        rand0.should.be.bignumber.not.equal(rand2);
        rand0.should.be.bignumber.equal('64856073772839990506814373782217928521534618466099710722049665631602958392435');
    });

    it('sum of computed skills is 250', async () => {
        for (let i = 0; i < 10; i++) {
            const seed = await assets.intHash("Barca" + i).should.be.fulfilled;
            const skills = await assets.computeSkills(seed).should.be.fulfilled;
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            sum.should.be.equal(250);
        }
    });

    it('get player pos in team', async () => {
        const nPLayersPerTeam = await assets.PLAYERS_PER_TEAM().should.be.fulfilled;
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const playerState = await assets.getPlayerState(playerId).should.be.fulfilled;
            const pos = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
            pos.toNumber().should.be.equal(playerId - 1);
        }
        await assets.getPlayerState(nPLayersPerTeam+1).should.be.rejected;
    })

    it('get existing virtual player skills', async () => {
        const numSkills = await assets.NUM_SKILLS().should.be.fulfilled;
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
        const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
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
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
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
        await assets.setPlayerState(state).should.be.fulfilled;
        const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
        const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
        skills[0].should.be.bignumber.equal('1');
        skills[1].should.be.bignumber.equal('2');
        skills[2].should.be.bignumber.equal('3');
        skills[3].should.be.bignumber.equal('4');
        skills[4].should.be.bignumber.equal('5');
    });

    it('compute player birth', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        const birth = await assets.computeBirth(0).should.be.fulfilled;
        birth.should.be.bignumber.equal('406');
    });

    it('get non virtual player team', async () => {
        await assets.createTeam("Barca",accounts[1]).should.be.fulfilled;
        await assets.createTeam("Madrid",accounts[1]).should.be.fulfilled;
        let playerState = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
        const teamBefore = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
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
        await assets.setPlayerState(state).should.be.fulfilled;
        playerState = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
        const teamAfter = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
        teamAfter.should.be.bignumber.not.equal(teamBefore);
        teamAfter.should.be.bignumber.equal('2');
    });

    it('create team', async () => {
        const receipt = await assets.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        let teamName = receipt.logs[0].args.teamName;
        teamName.should.be.equal("Barca",accounts[1]);
        const teamId = receipt.logs[0].args.teamId.toNumber();
        teamId.should.be.equal(1);
        teamName = await assets.getTeamName(teamId).should.be.fulfilled;
        teamName.should.be.equal("Barca",accounts[1]);
    });

    it('get playersId from teamId and pos in team', async () => {
        await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.rejected;
        await assets.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=11).should.be.rejected;
        let playerId = await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.fulfilled;
        playerId.toNumber().should.be.equal(1);
        playerId = await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=10).should.be.fulfilled;
        playerId.toNumber().should.be.equal(11);
    });

    it('sign team to league', async () => {
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
        await assets.createTeam(name = "Barca",accounts[1]).should.be.fulfilled;
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.fulfilled;
        const currentHistory = await assets.getTeamCurrentHistory(1).should.be.fulfilled;
        currentHistory.currentLeagueId.should.be.bignumber.equal('1');
        currentHistory.posInCurrentLeague.should.be.bignumber.equal('3');
        currentHistory.prevLeagueId.should.be.bignumber.equal('0');
        currentHistory.posInPrevLeague.should.be.bignumber.equal('0');
    });
});