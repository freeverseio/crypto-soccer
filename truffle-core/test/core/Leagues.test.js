require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;
    const blockInit = 1000;
    const blockStep = 10;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('default init block is 0', async () =>{
        const block = await leagues.getBlockInit().should.be.fulfilled;
        block.toNumber().should.be.equal(0);
    });

    it('default block step is 0', async () =>{
        const step = await leagues.getBlockStep().should.be.fulfilled;
        step.toNumber().should.be.equal(0);
    });

    it('default team ids is empty', async () => {
        const teamIds = await leagues.getTeamIds().should.be.fulfilled;
        teamIds.length.should.be.equal(0);
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
    })
});