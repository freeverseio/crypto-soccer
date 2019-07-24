const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Assets = artifacts.require('AssetsMock');
const PlayerStateLib = artifacts.require('PlayerState');
const PLAYERS_PER_TEAM = 25;
const initBlock = 1;
const step = 1;

contract('Assets', (accounts) => {
    let assets = null;
    let playerStateLib = null;
    let deployBlock = null;
    const ALICE = accounts[1];
    const BOB = accounts[2];
    
    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        await assets.setStatesContract(playerStateLib.address);
        deployBlock = await web3.eth.getBlockNumber().should.be.fulfilled;
    });
    
    it('create league', async () => {
        await assets.createLeague(nTeams = 2, deployBlock + 10, step).should.be.fulfilled;
        await assets.createLeague(nTeams = 3, deployBlock + 10, step).should.be.rejected; // only even num teams allowed
        await assets.createLeague(nTeams = 2, deployBlock - 10, step).should.be.rejected; // only init in future
    });
    return;    
    
    it('generate virtual player state', async () => {
        await assets.generateVirtualPlayerState(0).should.be.rejected;
        await assets.generateVirtualPlayerState(1).should.be.rejected;
        await assets.createTeam(name = "Barca", ALICE).should.be.fulfilled;
        await assets.generateVirtualPlayerState(1).should.be.fulfilled;
        await assets.generateVirtualPlayerState(PLAYERS_PER_TEAM+1).should.be.rejected;
    });
    
    it('compute seed', async () => {
        const seed = await assets.computeSeed(web3.utils.keccak256("ciao"), 55).should.be.fulfilled;
        seed.should.be.bignumber.equal('34043593120303183903741292954315585295064490957430510451016950480665910064957');
    });

    it('get team creation timestamp', async () => {
        await assets.getTeamCreationTimestamp(1).should.be.rejected;
        const receipt = await assets.createTeam(name = "Barca", ALICE).should.be.fulfilled;
        const blockNumber = receipt.receipt.blockNumber;
        const block = await web3.eth.getBlock(blockNumber).should.be.fulfilled;
        const timestamp = await assets.getTeamCreationTimestamp(1).should.be.fulfilled;
        timestamp.should.be.bignumber.equal(block.timestamp.toString());
    });

    it('player birth is generated from team creation timestamp', async () => {
        const name = "Barca"
        await assets.createTeam(name, ALICE).should.be.fulfilled;
        const teamCreationTimestamp = await assets.getTeamCreationTimestamp(1).should.be.fulfilled;
        const playerState = await assets.getPlayerState(5).should.be.fulfilled;
        const playerBirth = await playerStateLib.getMonthOfBirthInUnixTime(playerState).should.be.fulfilled;
        const posInTeam = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
        const nameHash = web3.utils.keccak256(web3.eth.abi.encodeParameter('string', name))
        const seed = await assets.computeSeed(nameHash, posInTeam).should.be.fulfilled;
        const computedBirth = await assets.computeBirth(seed, teamCreationTimestamp).should.be.fulfilled;
        playerBirth.should.be.bignumber.equal(computedBirth);
    });

    it('get playerIds of the team', async () => {
        await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
        let playerIds = await assets.getTeamPlayerIds(1).should.be.fulfilled;
        playerIds.length.should.be.equal(PLAYERS_PER_TEAM);
        for (let pos = 0; pos < PLAYERS_PER_TEAM ; pos++) 
            playerIds[pos].should.be.bignumber.equal((pos+1).toString());

        await assets.createTeam(name = "Madrid",ALICE).should.be.fulfilled;
        await assets.exchangePlayersTeams(playerId0 = PLAYERS_PER_TEAM, playerId1 = PLAYERS_PER_TEAM+3).should.be.fulfilled;
        playerIds = await assets.getTeamPlayerIds(1).should.be.fulfilled;
        playerIds.length.should.be.equal(PLAYERS_PER_TEAM);
        for (let pos = 0; pos < PLAYERS_PER_TEAM-1 ; pos++) 
            playerIds[pos].should.be.bignumber.equal((pos+1).toString());
        playerIds[PLAYERS_PER_TEAM-1].should.be.bignumber.equal((PLAYERS_PER_TEAM+3).toString());
    });

    it('add team with different owner than the sender', async () => {
        await assets.createTeam('Barca', ALICE).should.be.fulfilled;
        const owner = await assets.getTeamOwner(teamId = 1).should.be.fulfilled;
        owner.should.be.equal(ALICE);
    })

    it('add 2 teams with same name', async() => {
        await assets.createTeam('Barca', ALICE).should.be.fulfilled;
        await assets.createTeam('Barca', ALICE).should.be.rejected;
    })

    it('team exists', async () => {
        let result = await assets.teamExists(0).should.be.fulfilled;
        result.should.be.equal(false);
        result = await assets.teamExists(1).should.be.fulfilled;
        result.should.be.equal(false);
        await assets.createTeam("Barca", ALICE).should.be.fulfilled;
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
        await assets.getTeamNameHash(teamId = 0).should.be.rejected;
    });

    it('get name of unexistent team', async () => {
        await assets.getTeamNameHash(teamId = 1).should.be.rejected;
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
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
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
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        await assets.isVirtual(1).should.eventually.equal(true);
    });

    it('set player state of existent virtual player', async () => {
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        let state = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
        const currentBlock = 5; // TODO: get it properly
        state = await playerStateLib.setLastSaleBlock(state, currentBlock).should.be.fulfilled;
        await assets.setPlayerState(state).should.be.fulfilled;
        const resultState = await assets.getPlayerState(playerId).should.be.fulfilled;
        resultState.should.be.bignumber.equal(state);
    });

    it('is existent non virtual player', async () => {
        await assets.setPlayerState(4).should.be.rejected;
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
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

    it('get state of player on creation', async () => {
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        const state = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
        let result = await playerStateLib.getSkills(state).should.be.fulfilled;
        result.should.be.bignumber.equal('3819232821366079540');
        result = await playerStateLib.getPlayerId(state).should.be.fulfilled;
        result.should.be.bignumber.equal('1');
        result = await playerStateLib.getCurrentTeamId(state).should.be.fulfilled;
        result.should.be.bignumber.equal('1');
        result = await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
        result = await playerStateLib.getPrevLeagueId(state).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
        result = await playerStateLib.getPrevTeamPosInLeague(state).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
        result = await playerStateLib.getPrevShirtNumInLeague(state).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
        result = await playerStateLib.getLastSaleBlock(state).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
    });

    it('exchange players team', async () => {
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        await assets.createTeam("Madrid",ALICE).should.be.fulfilled;
        await assets.exchangePlayersTeams(playerId0 = 8, playerId1 = PLAYERS_PER_TEAM+3).should.be.fulfilled;
        const statePlayer0 = await assets.getPlayerState(playerId0).should.be.fulfilled;
        const teamPlayer0 = await playerStateLib.getCurrentTeamId(statePlayer0).should.be.fulfilled;
        teamPlayer0.should.be.bignumber.equal('2');
        const statePlayer1 = await assets.getPlayerState(playerId1).should.be.fulfilled;
        const teamPlayer1 = await playerStateLib.getCurrentTeamId(statePlayer1).should.be.fulfilled;
        teamPlayer1.should.be.bignumber.equal('1');
    });

    it('get player state of existing player', async () => {
        const nPLayersPerTeam = await assets.PLAYERS_PER_TEAM().should.be.fulfilled;
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
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
        rand0.should.be.bignumber.equal('16868380996023217686301278465084779672212597498847303814512224087959838246889');
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
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const playerState = await assets.getPlayerState(playerId).should.be.fulfilled;
            const pos = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
            pos.toNumber().should.be.equal(playerId - 1);
        }
        await assets.getPlayerState(nPLayersPerTeam+1).should.be.rejected;
    })

    it('get existing virtual player skills', async () => {
        const numSkills = await assets.NUM_SKILLS().should.be.fulfilled;
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
        const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
        skills.length.should.be.equal(numSkills.toNumber());
        skills[0].should.be.bignumber.equal('72');
        skills[1].should.be.bignumber.equal('71');
        skills[2].should.be.bignumber.equal('27');
        skills[3].should.be.bignumber.equal('42');
        skills[4].should.be.bignumber.equal('38');
        const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
        sum.should.be.equal(250);
    });


    it('get existing non virtual player skills', async () => {
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
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
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        const birth = await assets.computeBirth(0, 1557495456).should.be.fulfilled;
        birth.should.be.bignumber.equal('406');
    });

    it('get non virtual player team', async () => {
        await assets.createTeam("Barca",ALICE).should.be.fulfilled;
        await assets.createTeam("Madrid",ALICE).should.be.fulfilled;
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
        const receipt = await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        const teamId = receipt.logs[0].args.id.toNumber();
        teamId.should.be.equal(1);
        teamName = await assets.getTeamNameHash(teamId).should.be.fulfilled;
        expected = web3.utils.keccak256(web3.eth.abi.encodeParameter('string', name))
        teamName.should.be.equal(expected);
    });

    it('get playersId from teamId and pos in team', async () => {
        await assets.generateVirtualPlayerId(teamId = 1, posInTeam=0).should.be.rejected;
        await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
        await assets.generateVirtualPlayerId(teamId = 1, posInTeam=PLAYERS_PER_TEAM).should.be.rejected;
        let playerId = await assets.generateVirtualPlayerId(teamId = 1, posInTeam=0).should.be.fulfilled;
        playerId.toNumber().should.be.equal(1);
        playerId = await assets.generateVirtualPlayerId(teamId = 1, posInTeam=PLAYERS_PER_TEAM-1).should.be.fulfilled;
        playerId.toNumber().should.be.equal(PLAYERS_PER_TEAM);

    });

    it('sign team to league', async () => {
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
        await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.fulfilled;
        const currentHistory = await assets.getTeamCurrentHistory(1).should.be.fulfilled;
        currentHistory.currentLeagueId.should.be.bignumber.equal('1');
        currentHistory.posInCurrentLeague.should.be.bignumber.equal('3');
        currentHistory.prevLeagueId.should.be.bignumber.equal('0');
        currentHistory.posInPrevLeague.should.be.bignumber.equal('0');
    });

    it('sign team to league twice should fail', async () => {
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.rejected;
    });
    
    it('transfer team', async () => {
        await assets.createTeam(name = "Barca", ALICE).should.be.fulfilled;
        const currentOwner = await assets.getTeamOwner(teamId = 1).should.be.fulfilled;
        currentOwner.should.be.equal(ALICE);
        await assets.transferTeam(teamId = 1, BOB).should.be.fulfilled;
        const newOwner = await assets.getTeamOwner(teamId).should.be.fulfilled;
        newOwner.should.be.equal(BOB);
    });
        
    it('transfer non-exisiting team', async () => {
        await assets.transferTeam(teamId = 1, BOB).should.be.rejected;
    });

    it('transfer team accross same owner', async () => {
        await assets.createTeam(name = "Barca", ALICE).should.be.fulfilled;
        await assets.transferTeam(teamId = 1, ALICE).should.be.rejected;
    });
        
   
});