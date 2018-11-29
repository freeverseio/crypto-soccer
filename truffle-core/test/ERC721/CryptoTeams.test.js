require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('no initial teams', async () => {
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(0);
    })

    it('team owner', async () => {
        await contract.ownerOf(0).should.be.rejected;
        await contract.addTeam("team", accounts[0]).should.be.fulfilled;
        const owner = await contract.ownerOf(1).should.be.fulfilled;
        owner.should.be.equal(accounts[0]);
    })

    it('team name', async () => {
        const team = "team";
        await contract.getTeamName(0).should.be.rejected;
        await contract.getTeamName(1).should.be.rejected;
        await contract.addTeam(team, accounts[0]).should.be.fulfilled;
        const name = await contract.getTeamName(1).should.be.fulfilled;
        name.should.be.equal(team);
    });

    it('create team', async () => {
        await contract.addTeam("team", accounts[0]).should.be.fulfilled;
        const name = await contract.getTeamName(1).should.be.fulfilled;
        name.should.be.equal("team");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(1);
    })
});
