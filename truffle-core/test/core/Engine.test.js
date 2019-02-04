require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let leagues = null;
    let engine = null;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        engine = await Engine.new(leagues.address).should.be.fulfilled;
    });

    it('Leagues contract', async () => {
        const address = await engine.getLeaguesContract().should.be.fulfilled;
        address.should.be.equal(leagues.address);
    });
});