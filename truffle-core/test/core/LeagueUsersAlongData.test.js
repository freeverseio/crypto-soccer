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
    const initBlock = 1;
    const step = 1;
    const leagueId = 1;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 1;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        // engine = await Engine.new().should.be.fulfilled;
        leagues = await League.new().should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
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
        await leagues.updateUsersAlongDataHash(thisLeagueId = 2, teamIds = [0], tactic = [[4, 4, 2]], block = [3]).should.be.rejected;
    })

    it('update finished league', async () => {
        const finished = await leagues.hasFinished(leagueId).should.be.fulfilled;
        finished.should.be.equal(true);
        await leagues.updateUsersAlongDataHash(leagueId, teamIds = [0], tactic = [[4, 4, 2]], block = [3]).should.be.rejected;
    });
    
    return; // TODO: continue only when doing lionel4 properly
    
    it('update', async () => {
        await leagues.create(
            leagueId = 0, 
            initBlock = 1, 
            step = 100000, 
            teamIds = [1, 2], 
            tactics = [[4, 4, 3], [4, 4, 3]]
        ).should.be.fulfilled;
        const finished = await leagues.hasFinished(leagueId).should.be.fulfilled;
        finished.should.be.equal(false);
        await leagues.updateUsersAlongDataHash(leagueId, teamIds = [0], tactic = [[4, 4, 2]], block = [3]).should.be.fulfilled;
        const hash = await leagues.computeUsersAlongDataHash(teamIds = [0], tactic = [[4, 4, 2]], block = [3]).should.be.fulfilled;
        const usersAlongDataHash = await leagues.getUsersAlongDataHash(leagueId).should.be.fulfilled;
        hash.should.be.equal(usersAlongDataHash);
    });

    it('compute user along data hash', async () => {
        let hash = await leagues.computeUsersAlongDataHash(teamIds = [0], tactic = [[4, 4, 2]], block = [3]).should.be.fulfilled;
        hash.should.be.equal('0x23f31280f69accf85f4ed1f35b9b7c8120241435f7f1c7005d1a397e09035c4b');
        hash = await leagues.computeUsersAlongDataHash(teamIds = [0], tactic = [[4, 4, 2]], block = [2]).should.be.fulfilled;
        hash.should.be.equal('0x94cc21c8dfb0a81fb883059124ef97d417f42f86c1caa0c248ae05eda99ff245');
        hash = await leagues.computeUsersAlongDataHash(teamIds = [0, 1], tactic = [[4, 4, 2], [4, 4, 2]], block = [2, 4]).should.be.fulfilled;
        hash.should.be.equal('0xb854db7f540d4de46dd8e42fdaf48fed19057ae2e0e60e8be3c460647ceae2d6');
    });


    // it('update with wrong teamIdx', async () => {
    //     await leagues.create(
    //         leagueId = 0, 
    //         initBlock = 1, 
    //         step = 100000, 
    //         teamIds = [1, 2], 
    //         tactics = [[4, 4, 3], [4, 4, 3]]
    //     ).should.be.fulfilled;
    //     const finished = await leagues.hasFinished(leagueId).should.be.fulfilled;
    //     finished.should.be.equal(false);
    //     await leagues.updateUsersAlongDataHash(leagueId, teamIdx = [1, 3], tactic = [[4, 4, 2], [4, 4, 2]]).should.be.rejected;
    // });
}) 