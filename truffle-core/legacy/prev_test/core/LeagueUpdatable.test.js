require('chai')
    .use(require('chai-as-promised'))
    .should();


const Assets = artifacts.require('Assets');
const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');
const State = artifacts.require('LeagueState');
const GameController = artifacts.require('GameControllerDummy');

contract('LeaguesUpdatable', (accounts) => {
    let leagues = null;
    let gameControllerDummy = null;
    const leagueId = 1;
    let assets = null;
    let playerStateLib = null;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 1;
    const initBlock = 1;
    const step = 1;

    beforeEach(async () => {
        state = await State.new().should.be.fulfilled;
        assets = await Assets.new(state.address).should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, state.address).should.be.fulfilled;
        gameControllerDummy = await GameController.new().should.be.fulfilled;
        await leagues.setStakersContract(gameControllerDummy.address).should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 1, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 2, order, tactic442).should.be.fulfilled;
    });

    it('unexistent league', async () => {
        await leagues.getDayStateHashes(thhisLeagueId = 3).should.be.rejected;
        await leagues.getInitStateHash(thhisLeagueId).should.be.rejected;
    })

    it('default hashes values on create league', async () => {
        const initHash = await leagues.getInitStateHash(leagueId).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await leagues.getDayStateHashes(leagueId).should.be.fulfilled;
        finalHashes.length.should.be.equal(0);
    });

    it('is updated', async () => {
        let result = await leagues.isUpdated(leagueId).should.be.fulfilled;
        result.should.be.equal(false);
        const initStateHash = '0x54564';
        const dayStateHashes = ['0x24353', '0x5434432'];
        const scores = ['0x12', '0x3'];
        await leagues.updateLeague(leagueId, initStateHash, dayStateHashes, scores, isLie = false).should.be.fulfilled;
        result = await leagues.isUpdated(leagueId).should.be.fulfilled;
        result.should.be.equal(true);
    });
 
    it('updateBlock and updater', async () => {
        const initStateHash = '0x54564';
        const dayStateHashes = ['0x24353', '0x5434432'];
        const scores = ['0x12', '0x3'];
        const result = await leagues.updateLeague(leagueId, initStateHash, dayStateHashes, scores, isLie = false).should.be.fulfilled;
        const updateBlock = await leagues.getUpdateBlock(leagueId).should.be.fulfilled;
        updateBlock.toNumber().should.be.equal(result.receipt.blockNumber);
        const updater = await leagues.getUpdater(leagueId).should.be.fulfilled;
        updater.should.be.equal(accounts[0]);
    });
})