require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('add team', async () => {
        let count = await contract.totalSupply().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
        const id = 1; // TODO this smells : how I get this id from ? 
                      // I have to know the internal behavior of the component ---> WRONG
        await contract.addTeam("team", accounts[0]).should.be.fulfilled;
        count = await contract.totalSupply().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
    })
});
