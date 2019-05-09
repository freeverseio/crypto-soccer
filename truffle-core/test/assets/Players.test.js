require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('PlayersMock');

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

    it('get player team of existing player', async () => {
        await players.addTeam("Barca").should.be.fulfilled;
        await players.getPlayerTeam(1).should.be.fulfilled;
        await players.getPlayerTeam(2).should.be.fulfilled;
        return;
        await players.getPlayerTeam(3).should.be.fulfilled;
        await players.getPlayerTeam(4).should.be.fulfilled;
        await players.getPlayerTeam(5).should.be.fulfilled;
        await players.getPlayerTeam(6).should.be.fulfilled;
        await players.getPlayerTeam(7).should.be.fulfilled;
        await players.getPlayerTeam(8).should.be.fulfilled;
        await players.getPlayerTeam(9).should.be.fulfilled;
        await players.getPlayerTeam(10).should.be.fulfilled;
    })
});
 