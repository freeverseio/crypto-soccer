require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');
const CryptoPlayers = artifacts.require('CryptoPlayersTeamed');

contract('CryptoPlayersTeamed', (accounts) => {
    it('deployment', async () => {
        const cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
        await CryptoPlayers.new(cryptoTeams.address).should.be.fulfilled;
    });

    it('get team unexistent player', async () => {
        const cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
        const contract = await CryptoPlayers.new(cryptoTeams.address).should.be.fulfilled;

        const id = 1;
        await contract.getTeam(id).should.be.rejected;
    })

    // it('mint', async () => {
    //     const contract = await CryptoPlayers.new().should.be.fulfilled;
    //     const teamId = 1;
    //     const teamName = "panzerotto";
    //     await contract.mint(accounts[0], teamId);
    //     const name = await contract.getName(teamId).should.be.fulfilled;
    //     name.should.be.equal(name);
    // })
});
