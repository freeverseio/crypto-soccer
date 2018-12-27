require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeamsMetadataMock');

contract('CryptoTeamsMetadata', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('symbol', async () => {
        const symbol = await contract.symbol().should.be.fulfilled;
        symbol.should.be.equal("CST");
    });
    
    it('name', async () => {
        const name = await contract.name().should.be.fulfilled;
        name.should.be.equal("CryptoSoccerTeams");
    });

    it('tokenURI of unexistend team', async () => {
        await contract.tokenURI(0).should.be.rejected;
        await contract.tokenURI(1).should.be.rejected;
    });

    it('tokenURI of existent team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const uri = await contract.tokenURI(id).should.be.fulfilled;
        uri.should.be.equal("?playersId=0");
    });

    it('set URI without being URIer', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        await contract.renounceURIer().should.be.fulfilled;
        await contract.setTokensURI("URI").should.be.rejected;
    });
});
