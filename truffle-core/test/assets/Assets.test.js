const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Assets = artifacts.require('AssetsMock');
const PlayerStateLib = artifacts.require('PlayerState');
const PLAYERS_PER_TEAM = 25;
const TEAMS_PER_LEAGUE = 10;
const step = 1;

contract('Assets', (accounts) => {
    let assets = null;
    let playerStateLib = null;
    let deployBlock = null;
    let futureBlock = null;
    const ALICE = accounts[1];
    const BOB = accounts[2];
    
    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        await assets.setStatesContract(playerStateLib.address);
        deployBlock = await web3.eth.getBlockNumber().should.be.fulfilled;
        futureBlock = deployBlock + 200
    });
    
    // it('create league', async () => {
    //     let nLeagues = await assets.leaguesCount().should.be.fulfilled;
    //     nLeagues.toNumber().should.be.equal(0);
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     nLeagues = await assets.leaguesCount().should.be.fulfilled;
    //     nLeagues.toNumber().should.be.equal(1);
    // });

    // it('create league with wrong init params', async () => {
    //     await assets.createLeague(deployBlock - 10, step).should.be.rejected; // only init in future
    //     await assets.createLeague(deployBlock - 10, thisStep = -20).should.be.rejected; // only positive
    // });

    // it('playerExists upon league creation', async () => {
    //     let exists = await assets.playerExists(0).should.be.fulfilled; // playerId = 0 is dummy
    //     exists.should.be.equal(false);
    //     exists = await assets.playerExists(1).should.be.fulfilled;
    //     exists.should.be.equal(false);
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     exists = await assets.playerExists(0).should.be.fulfilled;
    //     exists.should.be.equal(false);
    //     exists = await assets.playerExists(1).should.be.fulfilled;
    //     exists.should.be.equal(true);
    //     exists = await assets.playerExists(PLAYERS_PER_TEAM).should.be.fulfilled;
    //     exists.should.be.equal(true);
    //     exists = await assets.playerExists(TEAMS_PER_LEAGUE*PLAYERS_PER_TEAM).should.be.fulfilled;
    //     exists.should.be.equal(true);
    //     exists = await assets.playerExists(TEAMS_PER_LEAGUE*PLAYERS_PER_TEAM+1).should.be.fulfilled;
    //     exists.should.be.equal(false);
    // });
    
    // it('generate virtual player state', async () => {
    //     await assets.generateVirtualPlayerState(0).should.be.rejected; // playerId = 0 is dummy
    //     await assets.generateVirtualPlayerState(1).should.be.rejected;
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     await assets.generateVirtualPlayerState(0).should.be.rejected; 
    //     await assets.generateVirtualPlayerState(1).should.be.fulfilled;
    //     await assets.generateVirtualPlayerState(TEAMS_PER_LEAGUE*PLAYERS_PER_TEAM).should.be.fulfilled;
    //     await assets.generateVirtualPlayerState(TEAMS_PER_LEAGUE*PLAYERS_PER_TEAM+1).should.be.rejected;
    // });
    
    // it('compute seed', async () => {
    //     const seed = await assets.computeSeed(web3.utils.keccak256("ciao"), 55).should.be.fulfilled;
    //     seed.should.be.bignumber.equal('34043593120303183903741292954315585295064490957430510451016950480665910064957');
    // });

    // it('get team creation timestamp', async () => {
    //     await assets.getTeamCreationBlocknum(teamId = 1).should.be.rejected;
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const blockReported = await assets.getTeamCreationBlocknum(teamId).should.be.fulfilled;
    //     blockReported.toNumber().should.be.equal(futureBlock);
    // });

    // it('player birth is generated from team creation timestamp', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     teamId = 3;
    //     posInTeam = 4;
    //     playerId = (teamId-1)*PLAYERS_PER_TEAM+ posInTeam +1
    //     leagueId = 1;
    //     posInLeague = teamId-1; // pos = 0 is for first team, so this is pos = 1 

    //     const teamCreationBlocknum = await assets.getTeamCreationBlocknum(teamId).should.be.fulfilled;
    //     teamCreationBlocknum.toNumber().should.be.equal(futureBlock);

    //     const playerState = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     const playerBirth = await playerStateLib.getMonthOfBirthInUnixTime(playerState).should.be.fulfilled;
    //     const posInTeamReported = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
    //     posInTeamReported.toNumber().should.be.equal(posInTeam)

    //     const dna = web3.utils.keccak256(web3.eth.abi.encodeParameters(
    //             ['uint256', 'uint8'],
    //             [leagueId, posInLeague]
    //         )
    //     );
    //     const dnaReported = await assets.botTeamIdToDNA(teamId).should.be.fulfilled;
    //     dnaReported.should.be.equal(dna);

    //     const seed = await assets.computeSeed(dna, posInTeam).should.be.fulfilled;
    //     const computedBirth = await assets.computeBirth(seed, futureBlock).should.be.fulfilled;
    //     playerBirth.should.be.bignumber.equal(computedBirth);
    // });

    // it('get playerIds of the team', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     let playerIds = await assets.getTeamPlayerIds(teamId = 1).should.be.fulfilled;
    //     playerIds.length.should.be.equal(PLAYERS_PER_TEAM);
    //     for (let pos = 0; pos < PLAYERS_PER_TEAM ; pos++) 
    //         playerIds[pos].should.be.bignumber.equal((pos+1).toString());

    // });

    // it('get playerIds of the team after player transfer', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     let playerIds = await assets.getTeamPlayerIds(teamId = 1).should.be.fulfilled;
    //     playerIds.length.should.be.equal(PLAYERS_PER_TEAM);
    //     for (let pos = 0; pos < PLAYERS_PER_TEAM ; pos++) 
    //         playerIds[pos].should.be.bignumber.equal((pos+1).toString());
    //     // Players in bot teams cannot be traded:
    //     await assets.exchangePlayersTeams(playerId0 = PLAYERS_PER_TEAM, playerId1 = PLAYERS_PER_TEAM+3).should.be.rejected;
    //     await assets.transferTeam(teamId=1,ALICE);
    //     await assets.transferTeam(teamId=2,BOB);
    //     await assets.exchangePlayersTeams(playerId0 = PLAYERS_PER_TEAM, playerId1 = PLAYERS_PER_TEAM+3).should.be.fulfilled;
    //     playerIds = await assets.getTeamPlayerIds(teamId = 1).should.be.fulfilled;
    //     playerIds.length.should.be.equal(PLAYERS_PER_TEAM);
    //     for (let pos = 0; pos < PLAYERS_PER_TEAM-1 ; pos++) 
    //         playerIds[pos].should.be.bignumber.equal((pos+1).toString());
    //     playerIds[PLAYERS_PER_TEAM-1].should.be.bignumber.equal(playerId1.toString());
    // });

    // it('transfer a bot team', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     let isBot = await assets.isBotTeam(teamId = 1).should.be.fulfilled;
    //     isBot.should.be.equal(true);
    //     await assets.transferTeam(teamId, ALICE);
    //     isBot = await assets.isBotTeam(teamId).should.be.fulfilled;
    //     isBot.should.be.equal(false);
    //     let owner = await assets.getTeamOwner(teamId).should.be.fulfilled;
    //     owner.should.be.equal(ALICE);
    // });
    
    // it('team exists', async () => {
    //     let result = await assets.teamExists(0).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     result = await assets.teamExists(1).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     result = await assets.teamExists(1).should.be.fulfilled;
    //     result.should.be.equal(true);
    //     result = await assets.teamExists(TEAMS_PER_LEAGUE + 1).should.be.fulfilled;
    //     result.should.be.equal(false);
    // });

    // it('initial number of team', async () => {
    //     const count = await assets.countTeams().should.be.fulfilled;
    //     count.toNumber().should.be.equal(0);
    // });

    // it('get dna of invalid team', async () => {
    //     await assets.getTeamDNA(teamId = 0).should.be.rejected;
    // });

    // it('get name of unexistent team', async () => {
    //     await assets.getTeamDNA(teamId = 1).should.be.rejected;
    // });

    // it('existence of null player', async () => {
    //     const exists = await assets.playerExists(playerId = 0).should.be.fulfilled;
    //     exists.should.be.equal(false);
    // });

    // it('existence of unexistent player', async () => {
    //     const exists = await assets.playerExists(playerId = 1).should.be.fulfilled;
    //     exists.should.be.equal(false);
    // });

    // it('existence of existent player', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const exists = await assets.playerExists(playerId = 1).should.be.fulfilled;
    //     exists.should.be.equal(true);
    // });

    // it('is null player virtual', async () => {
    //     await assets.isPlayerVirtual(0).should.be.rejected;
    // });

    // it('is unexistent player virtual', async () => {
    //     await assets.isPlayerVirtual(1).should.be.rejected;
    // });

    // it('is existent player virtual', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     await assets.isPlayerVirtual(1).should.eventually.equal(true);
    // });
    
    // it('set player state of existent virtual player', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     let state = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
    //     const currentBlock = 5; // TODO: get it properly
    //     state = await playerStateLib.setLastSaleBlock(state, currentBlock).should.be.fulfilled;
    //     await assets.setPlayerState(state).should.be.fulfilled;
    //     const resultState = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     resultState.should.be.bignumber.equal(state);
    // });
    
    // it('is existent non virtual player', async () => {
    //     await assets.setPlayerState(4).should.be.rejected;
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const state = await playerStateLib.playerStateCreate(
    //         defence = 3,
    //         speed = 3,
    //         pass = 3,
    //         shoot = 3,
    //         endurance = 3,
    //         monthOfBirthInUnixTime = 3,
    //         playerId = 1,
    //         currentTeamId = 1,
    //         currentShirtNum = 3,
    //         prevLeagueId = 3,
    //         prevTeamPosInLeague = 3,
    //         prevShirtNumInLeague = 3,
    //         lastSaleBlock = 3
    //     ).should.be.fulfilled;
    //     await assets.setPlayerState(state).should.be.fulfilled;
    //     await assets.isPlayerVirtual(playerId = 1).should.eventually.equal(false);
    // });
    

    // it('get state of player on creation', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const state = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
    //     let result = await playerStateLib.getSkills(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('2522270891984633908');
    //     result = await playerStateLib.getPlayerId(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('1');
    //     result = await playerStateLib.getCurrentTeamId(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('1');
    //     result = await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('0');
    //     result = await playerStateLib.getPrevLeagueId(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('0');
    //     result = await playerStateLib.getPrevTeamPosInLeague(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('0');
    //     result = await playerStateLib.getPrevShirtNumInLeague(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('0');
    //     result = await playerStateLib.getLastSaleBlock(state).should.be.fulfilled;
    //     result.should.be.bignumber.equal('0');
    // });

    // it('get player state of existing player', async () => {
    //     const nPLayersPerTeam = await assets.PLAYERS_PER_TEAM().should.be.fulfilled;
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++)
    //         await assets.getPlayerState(playerId).should.be.fulfilled;
    //     await assets.getPlayerState(nPLayersPerTeam+1).should.be.rejected;
    // });

    // it('computed skills with rnd = 0 is 50 each', async () => {
    //     let skills = await assets.computeSkills(0).should.be.fulfilled;
    //     skills.forEach(skill => (skill.toNumber().should.be.equal(50)));
    // });

    // it('int hash is deterministic', async () => {
    //     const rand0 = await assets.intHash("Barca0").should.be.fulfilled;
    //     const rand1 = await assets.intHash("Barca0").should.be.fulfilled;
    //     rand0.should.be.bignumber.equal(rand1);
    //     const rand2 = await assets.intHash("Barca1").should.be.fulfilled;
    //     rand0.should.be.bignumber.not.equal(rand2);
    //     rand0.should.be.bignumber.equal('16868380996023217686301278465084779672212597498847303814512224087959838246889');
    // });

    // it('sum of computed skills is 250', async () => {
    //     for (let i = 0; i < 10; i++) {
    //         const seed = await assets.intHash("Barca" + i).should.be.fulfilled;
    //         const skills = await assets.computeSkills(seed).should.be.fulfilled;
    //         const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
    //         sum.should.be.equal(250);
    //     }
    // });

    // it('get player pos in team', async () => {
    //     const nPLayersPerTeam = await assets.PLAYERS_PER_TEAM().should.be.fulfilled;
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
    //         const playerState = await assets.getPlayerState(playerId).should.be.fulfilled;
    //         const pos = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
    //         pos.toNumber().should.be.equal(playerId - 1);
    //     }
    //     await assets.getPlayerState(nPLayersPerTeam+1).should.be.rejected;
    // })

    // it('get existing virtual player skills', async () => {
    //     const numSkills = await assets.NUM_SKILLS().should.be.fulfilled;
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
    //     const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
    //     skills.length.should.be.equal(numSkills.toNumber());
    //     skills[0].should.be.bignumber.equal('69');
    //     skills[1].should.be.bignumber.equal('55');
    //     skills[2].should.be.bignumber.equal('45');
    //     skills[3].should.be.bignumber.equal('25');
    //     skills[4].should.be.bignumber.equal('56');
    //     const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
    //     sum.should.be.equal(250);
    // });


    // it('get existing non virtual player skills', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const state = await playerStateLib.playerStateCreate(
    //         defence = 1,
    //         speed = 2,
    //         pass = 3,
    //         shoot = 4,
    //         endurance = 5,
    //         monthOfBirthInUnixTime = 6,
    //         playerId = 10,
    //         currentTeamId = 1,
    //         currentShirtNum = 3,
    //         prevLeagueId = 3,
    //         prevTeamPosInLeague = 3,
    //         prevShirtNumInLeague = 3,
    //         lastSaleBlock = 3
    //     ).should.be.fulfilled;
    //     await assets.setPlayerState(state).should.be.fulfilled;
    //     const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
    //     const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
    //     skills[0].should.be.bignumber.equal('1');
    //     skills[1].should.be.bignumber.equal('2');
    //     skills[2].should.be.bignumber.equal('3');
    //     skills[3].should.be.bignumber.equal('4');
    //     skills[4].should.be.bignumber.equal('5');
    // });

    // it('compute player birth', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const birth = await assets.computeBirth(0, 1557495456).should.be.fulfilled;
    //     birth.should.be.bignumber.equal('406');
    // });

    // it('get non virtual player team', async () => {
    //     await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     let playerState = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
    //     const teamBefore = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
    //     const state = await playerStateLib.playerStateCreate(
    //         defence = 3,
    //         speed = 3,
    //         pass = 3,
    //         shoot = 3,
    //         endurance = 3,
    //         monthOfBirthInUnixTime = 3,
    //         playerId = 1,
    //         currentTeamId = 2,
    //         currentShirtNum = 3,
    //         prevLeagueId = 3,
    //         prevTeamPosInLeague = 3,
    //         prevShirtNumInLeague = 3,
    //         lastSaleBlock = 3
    //     ).should.be.fulfilled;
    //     await assets.setPlayerState(state).should.be.fulfilled;
    //     playerState = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
    //     const teamAfter = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
    //     teamAfter.should.be.bignumber.not.equal(teamBefore);
    //     teamAfter.should.be.bignumber.equal('2');
    // });
    
    // it('creation of teams through league creation', async () => {
    //     const receipt = await assets.createLeague(futureBlock, step).should.be.fulfilled;
    //     const count = await assets.countTeams().should.be.fulfilled;
    //     count.toNumber().should.be.equal(TEAMS_PER_LEAGUE);
    //     const leagueId = receipt.logs[0].args.leagueId.toNumber();
    //     leagueId.should.be.equal(1);
    // });

    it('sign team to league', async () => {
        await assets.signToLeague(teamId = 5, leagueId = 1, posInLeague = 0).should.be.rejected;
        await assets.createLeague(futureBlock, step).should.be.fulfilled;
        await assets.createLeague(futureBlock + 10 * step, step).should.be.fulfilled;
        await assets.transferTeam(teamId, ALICE);
        currentHistory = await assets.getTeamCurrentHistory(teamId).should.be.fulfilled;
        currentHistory.currentLeagueId.should.be.bignumber.equal('1');
        currentHistory.posInCurrentLeague.should.be.bignumber.equal('4');
        currentHistory.prevLeagueId.should.be.bignumber.equal('0');
        currentHistory.posInPrevLeague.should.be.bignumber.equal('0');
        await assets.signToLeague(teamId, leagueId = 2, posInLeague = 3).should.be.fulfilled;
        currentHistory = await assets.getTeamCurrentHistory(teamId).should.be.fulfilled;
        currentHistory.currentLeagueId.should.be.bignumber.equal('2');
        currentHistory.posInCurrentLeague.should.be.bignumber.equal('3');
        currentHistory.prevLeagueId.should.be.bignumber.equal('1');
        currentHistory.posInPrevLeague.should.be.bignumber.equal('4');
    });

    it('sign team to league where it already belongs', async () => {
        await assets.createLeague(futureBlock, step).should.be.fulfilled;
        await assets.createLeague(futureBlock + step, step).should.be.fulfilled;
        await assets.transferTeam(teamId, ALICE);
        // it already belongs to league = 1:
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
        // await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.rejected;
        // it can only sign once to league = 2:
        // await assets.signToLeague(teamId = 1, leagueId = 2, posInLeague = 3).should.be.fulfilled;
        // await assets.signToLeague(teamId = 1, leagueId = 2, posInLeague = 3).should.be.rejected;
    });

    it('transfer non-bots team', async () => {
        await assets.createLeague(futureBlock, step).should.be.fulfilled;
        await assets.transferTeam(teamId = 1, ALICE);
        const currentOwner = await assets.getTeamOwner(teamId).should.be.fulfilled;
        currentOwner.should.be.equal(ALICE);
        await assets.transferTeam(teamId, BOB).should.be.fulfilled;
        const newOwner = await assets.getTeamOwner(teamId).should.be.fulfilled;
        newOwner.should.be.equal(BOB);
    });
        
    it('transfer non-exisiting team', async () => {
        await assets.transferTeam(teamId = 1, BOB).should.be.rejected;
    });

    it('transfer team accross same owner', async () => {
        await assets.createLeague(futureBlock, step).should.be.fulfilled;
        await assets.transferTeam(teamId = 1, ALICE);
        await assets.transferTeam(teamId = 1, ALICE).should.be.rejected;
    });
        
   
});