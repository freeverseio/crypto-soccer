require('chai')
    .use(require('chai-as-promised'))
    .should();

const Assets = artifacts.require('Teams');

contract('Assets', (accounts) => {
    let assets = null;

    beforeEach(async () => {
        assets = await Assets.new().should.be.fulfilled;
    });

    it('initial number of team', async () => {
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('create team', async () => {
        const receipt = await assets.createTeam(name = "Barca").should.be.fulfilled;
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        console.log(receipt)
        
    });
})