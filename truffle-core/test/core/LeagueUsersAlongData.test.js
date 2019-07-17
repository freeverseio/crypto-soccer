const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueUsersAlongData');
const Assets = artifacts.require('Assets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('LeagueUsersAlongData', (accounts) => {
    let leagues = null;
    let assets = null;
    let playerStateLib = null;
    let receipt = null;
    const initBlock = 1;
    const step = 1;
    const leagueId = 1;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 0;
    const tactic541 = 1;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        // engine = await Engine.new().should.be.fulfilled;
        leagues = await League.new().should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
        receipt = await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 1, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 2, order, tactic442).should.be.fulfilled;
    });

        
    it('initial hash of unexistent league', async () => {
        await leagues.getUsersAlongDataHash(thisLeagueId = 3).should.be.rejected;
    });
    
    it('initial hash of existing league', async () => {
        const hash = await leagues.getUsersAlongDataHash(leagueId).should.be.fulfilled;
        hash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
    });

    it('update unexistent league', async () => {
        await leagues.updateUsersAlongDataHash(thisLeagueId = 2, teamIds = [0], tactic = [tactic442], block = [3]).should.be.rejected;
    })

    it('update finished league', async () => {
        // first league, created at block = 1, has clearly finished, so it should not admit an update of tacticsIds
        let finished = await leagues.hasFinished(leagueId).should.be.fulfilled;
        finished.should.be.equal(true);
        await leagues.updateUsersAlongDataHash(leagueId, teamIds = [0], tactic = [tactic442], block = [3]).should.be.rejected;
    });

    
    it('update unfinished league', async () => {
        const currentBlock = receipt.receipt.blockNumber;
        await leagues.create(nTeams = 2, thisBlockInit = currentBlock-1, thisStep = 300).should.be.fulfilled;
        finished = await leagues.hasFinished(thisLeagueId = 2).should.be.fulfilled;
        finished.should.be.equal(false);
        await leagues.updateUsersAlongDataHash(leagueId, teamIds = [0], tactic = [tactic442], block = [3]).should.be.rejected;
        await leagues.updateUsersAlongDataHash(leagueId, teamIds = [0], tactic = [tactic442], block = [3]).should.be.rejected;
        await leagues.updateUsersAlongDataHash(thisLeagueId = 2, teamIds = [0], tactic = [tactic442], block = [3]).should.be.fulfilled;
    });

    it('compute user along data hash', async () => {
        let hash = await leagues.computeUsersAlongDataHash(teamIds = [0], tactic = [tactic442], block = [3]).should.be.fulfilled;
        hash.should.be.equal('0xe962042d5dd1b0f1e2739622bbf6d76c2642bff7a104ac8c986507c4d6c1115d');
        hash = await leagues.computeUsersAlongDataHash(teamIds = [0], tactic = [tactic442], block = [2]).should.be.fulfilled;
        hash.should.be.equal('0xff496e08a30e0406970dcb66bf9f6ada8a180c31ceba28262332b01aa1921c70');
        hash = await leagues.computeUsersAlongDataHash(teamIds = [0, 1], tactic = [tactic442, tactic442], block = [2, 4]).should.be.fulfilled;
        hash.should.be.equal('0xb2c41e13c3b1e9eab924b357e919d04aea5f7a2a57a3b6a74f1f4da700cc5000');
    });


    // it('update with wrong teamIdx', async () => {
    //     await leagues.create(
    //         leagueId = 0, 
    //         initBlock = 1, 
    //         step = 100000, 
    //         teamIds = [1, 2], 
    //         tacticsIds = [[4, 4, 3], [4, 4, 3]]
    //     ).should.be.fulfilled;
    //     const finished = await leagues.hasFinished(leagueId).should.be.fulfilled;
    //     finished.should.be.equal(false);
    //     await leagues.updateUsersAlongDataHash(leagueId, teamIdx = [1, 3], tactic = [tactic442, tactic442]).should.be.rejected;
    // });
}) 