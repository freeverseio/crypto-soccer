require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeamsMetadata');

contract('CryptoTeamsBase', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('tokenURI of unexistend team', async () => {
        await contract.tokenURI(0).should.be.rejected;
        await contract.tokenURI(1).should.be.rejected;
    });

    it('tokenURI of existent team', async () => {
        const id = 1; 
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        const uri = await contract.tokenURI(id).should.be.fulfilled;
        uri.should.be.equal("");
    });
});
