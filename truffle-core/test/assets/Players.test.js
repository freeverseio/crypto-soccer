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
        const nPLayersPerTeam = await players.getPlayersPerTeam().should.be.fulfilled;
        await players.addTeam("Barca").should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const teamId = await players.getPlayerTeam(playerId).should.be.fulfilled;
            teamId.toNumber().should.be.equal(1);
        }
        await players.getPlayerTeam(nPLayersPerTeam+1).should.be.rejected;
    })
});
 