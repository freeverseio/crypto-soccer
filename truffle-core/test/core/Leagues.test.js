require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('default init hash is 0', async () =>{
        const init = await leagues.getInit().should.be.fulfilled;
        init.toNumber().should.be.equal(0);
    });

    it('default final hash is 0', async () =>{
        const final = await leagues.getFinal().should.be.fulfilled;
        final.toNumber().should.be.equal(0);
    });
});