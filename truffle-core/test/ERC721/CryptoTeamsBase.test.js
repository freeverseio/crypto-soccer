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
        const id = 0;
        await contract.mintWithTeamName(accounts[0], id, "team").should.be.fulfilled;
    });

    it('mint team', async () => {
        const id = 1;
        await contract.mintWithTeamName(accounts[0], id, "team").should.be.fulfilled;
        const name = await contract.getTeamName(id).should.be.fulfilled;
        name.should.be.equal("team");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(id);
    });

    it('mint team with an existent id', async () =>  {
        const id = 1;
        await contract.mintWithTeamName(accounts[0], id, "team").should.be.fulfilled;
        await contract.mintWithTeamName(accounts[0], id, "team").should.be.rejected;
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
