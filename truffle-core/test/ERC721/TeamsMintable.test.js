require('chai')
    .use(require('chai-as-promised'))
    .should();

const Teams = artifacts.require('TeamsMintable');

contract('TeamsMintable', (accounts) => {
    let contract = null;
    let Players = null;

    beforeEach(async () => {
        contract = await Teams.new().should.be.fulfilled;
    });

    it('mint team', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(1);
    });

    it('mint team with same name is forbidden', async () =>  {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        await contract.mint(accounts[0], "team").should.be.rejected;
    });

    it('name of minted team', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
    });

    it('get team id', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
    });

    it('get team id of unexistent tema', async () => {
        const id = await contract.getTeamId("team").should.be.rejected;
    });
});