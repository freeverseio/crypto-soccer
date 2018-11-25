require('chai')
    .use(require('chai-as-promised'))
    .should();

const TeamFactory = artifacts.require('TeamFactory');
const GameEngine = artifacts.require('GameEngine');

contract('GameEngine', (accounts) => {
    let teamFactory;
    let instance;

    beforeEach(async () => {
        teamEngine = await TeamFactory.new().should.be.fulfilled;
        instance = await GameEngine.new(teamEngine.address).should.be.fulfilled;
    });

    it('check name and symbol', async () => {
        console.log(instance.address);
    });
});