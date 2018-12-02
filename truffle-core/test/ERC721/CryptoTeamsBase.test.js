require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeamsBase');

contract('CryptoTeamsBase', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('no initial teams', async () => {
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(0);
    })

    it('id 0 is forbidden', async () => {
        await contract.mintWithTeamName(accounts[0], 0, "team").should.be.fulfilled;
    });

    it('mint team', async () => {
        await contract.mintWithTeamName(accounts[0], 1, "team").should.be.fulfilled;
        const name = await contract.getTeamName(1).should.be.fulfilled;
        name.should.be.equal("team");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(1);
    });

    it('team owner', async () => {
        const id = 1;
        await contract.mintWithTeamName(accounts[0], id, "team").should.be.fulfilled;
        const owner = await contract.ownerOf(id).should.be.fulfilled;
        owner.should.be.equal(accounts[0]);
    });

    it('team name', async () => {
        const id = 1;
        await contract.mintWithTeamName(accounts[0], id, "team").should.be.fulfilled;
        const name = await contract.getTeamName(id).should.be.fulfilled;
        name.should.be.equal("team");
    });


});
