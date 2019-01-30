require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');

leagues('Leagues', (accounts) => {
    let leagues = null;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    })
});