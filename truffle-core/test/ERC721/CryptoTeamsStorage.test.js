require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeamsStorageMock');

contract('CryptoTeamsStorage', (accounts) => {
    let contract = null;
    let cryptoPlayers = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('no initial teams', async () => {
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(0);
    })

    it('mint team', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(1);
    });

    it('team name', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        await contract.setName(id, "team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
    });

    it('get name of unexistent team', async () => {
        await contract.getName(1).should.be.rejected;
    });

    it('new team has no players', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const players = await contract.getPlayers(id).should.be.fulfilled;
        players.length.should.be.equal(0);
    });
});