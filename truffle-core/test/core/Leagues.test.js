require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;
    const blockInit = 1000;
    const blockStep = 10;
    const id = 1;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('default values', async () =>{
        const block = await leagues.getBlockInit().should.be.fulfilled;
        block.toNumber().should.be.equal(0);
        const step = await leagues.getBlockStep().should.be.fulfilled;
        step.toNumber().should.be.equal(0);
        const teamIds = await leagues.getTeamIds().should.be.fulfilled;
        teamIds.length.should.be.equal(0);
        const initStateHash = await leagues.getInitStateHash(id).should.be.fulfilled;
        initStateHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const stateHash = await leagues.getStateHash(id).should.be.fulfilled;
        stateHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
    })

    it('create league with no team', async () => {
        const teamIds = [];
        await leagues.create(blockInit, blockStep, teamIds).should.be.rejected;
    });

    it('create league with 1 team', async () => {
        const teamIds = [1];
        await leagues.create(blockInit, blockStep, teamIds).should.be.rejected;
    });

    it('create league with 2 teams', async () => {
        const teamIds = [1, 2];
        await leagues.create(blockInit, blockStep, teamIds).should.be.fulfilled;
        const result = await leagues.getTeamIds().should.be.fulfilled;
        result.length.should.be.equal(2);
        result[0].toNumber().should.be.equal(1);
        result[1].toNumber().should.be.equal(2);
    });

    it('check if not initialized league has started', async () => {
        await leagues.hasStarted().should.be.rejected;
    });

    it('is future league started', async () => {
        const teamIds = [1, 2];
        await leagues.create(blockInit, blockStep, teamIds).should.be.fulfilled;
        const isStarted = await leagues.hasStarted().should.be.fulfilled;
        isStarted.should.be.equal(false);
    });

    it('is current league started', async () => {
        const teamIds = [1, 2];
        const blockInit = 1;
        await leagues.create(blockInit, blockStep, teamIds).should.be.fulfilled;
        const isStarted = await leagues.hasStarted().should.be.fulfilled;
        isStarted.should.be.equal(true);
    });
});