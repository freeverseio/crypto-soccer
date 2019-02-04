require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let leagues = null;
    let engine = null;
    const id = 0;
    const initPlayerState = 0;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        engine = await Engine.new(leagues.address).should.be.fulfilled;
    });

    it('Leagues contract', async () => {
        const address = await engine.getLeaguesContract().should.be.fulfilled;
        address.should.be.equal(leagues.address);
    });

    it('compute unexistent league', async () => {
        await engine.computeLeagueFinalState(id, initPlayerState).should.be.rejected;
    });

    it('compute league', async () => {
        const blocksToInit = 1;
        const step = 1;
        const teamIds = [1, 2];
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const scores = await engine.computeLeagueFinalState(id, initPlayerState).should.be.fulfilled;
        console.log(scores);
    });
});