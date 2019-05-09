require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');

contract('Players', (accounts) => {
    let players = null;

    beforeEach(async () => {
        players = await Players.new().should.be.fulfilled;
    });

    it('query null player id', async () => {
        await players.getPlayerTeam(0).should.be.rejected;
    });

    it('query non created player id', async () => {
        await players.getPlayerTeam(1).should.be.rejected;
    });
});
 