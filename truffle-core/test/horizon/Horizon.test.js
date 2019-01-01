require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');
const CryptoTeams = artifacts.require('CryptoTeams');
const Horizon = artifacts.require('Horizon');

contract('Horizon', (accounts) => {
    let instance = null;
    let cryptoPlayers = null;
    let cryptoTeams = null;

    beforeEach(async () => {
        cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
        cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
        instance = await Horizon.new(cryptoPlayers.address, cryptoTeams.address).should.be.fulfilled;
        await cryptoPlayers.addMinter(instance.address).should.be.fulfilled;
        await cryptoTeams.addMinter(instance.address).should.be.fulfilled;
        await cryptoPlayers.renounceMinter().should.be.fulfilled;
        await cryptoTeams.renounceMinter().should.be.fulfilled;
        await cryptoPlayers.setTeamsContract(cryptoTeams.address).should.be.fulfilled;
        await cryptoTeams.setPlayersContract(cryptoPlayers.address).should.be.fulfilled;
    });

    it('create Team', async () => {
        await instance.createTeam("team").should.be.fulfilled;
        const teamCount = await cryptoTeams.totalSupply().should.be.fulfilled;
        teamCount.toNumber().should.be.equal(1);
        const playerCount = await cryptoPlayers.totalSupply().should.be.fulfilled;
        playerCount.toNumber().should.be.equal(11);
    })
});