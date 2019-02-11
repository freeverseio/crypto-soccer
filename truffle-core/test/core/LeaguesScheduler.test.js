const keccak = require('keccak');
require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('LeaguesScheduler');

contract('LeaguesScheduler', (accounts) => {
    let leagues = null;
    const id = 0;

    beforeEach(async () => {
        const blocksToInit = 0;
        const step = 1
        const teamIds = [1, 2];
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
    });
});