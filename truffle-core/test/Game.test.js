require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');
const CryptoTeams = artifacts.require('CryptoTeams');
const Horizon = artifacts.require('Horizon');
const Leagues = artifacts.require('Leagues');

contract('Horizon', (accounts) => {
    let horizon = null;
    let cryptoPlayers = null;
    let cryptoTeams = null;
    let leagues = null;

    beforeEach(async () => {
        cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
        cryptoTeams = await CryptoTeams.new(cryptoPlayers.address).should.be.fulfilled;
        horizon = await Horizon.new(cryptoTeams.address).should.be.fulfilled;
        await cryptoPlayers.addMinter(horizon.address).should.be.fulfilled;
        await cryptoPlayers.renounceMinter().should.be.fulfilled;
        await cryptoTeams.addMinter(horizon.address).should.be.fulfilled;
        await cryptoTeams.renounceMinter().should.be.fulfilled;
        await cryptoPlayers.addTeamsContract(cryptoTeams.address).should.be.fulfilled;
        await cryptoPlayers.renounceTeamsContract().should.be.fulfilled;

        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('play a league of 2 teams', async () => {
        await horizon.createTeam("Barcelona").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        const barcelonaId = await cryptoTeams.getTeamId("Barcelona").should.be.fulfilled;
        const madridId = await cryptoTeams.getTeamId("Madrid").should.be.fulfilled;
        const blockInitDelta = 1;
        const step = 1;
        await leagues.create(blockInitDelta, step, [barcelonaId, madridId]).should.be.fulfilled;
    });
})
