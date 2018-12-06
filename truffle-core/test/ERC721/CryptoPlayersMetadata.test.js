require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMetadataMock');

contract('CryptoPlayersMetadata', (accounts) => {
    let contract = null;
    const URI = "QmUC4KA1Vi3DizRrTj9Z4uyrL6a7zjS7wNnvR5iNzYALSh";

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    })

   it('symbol', async () => {
        const symbol = await contract.symbol().should.be.fulfilled;
        symbol.should.be.equal("CSP");
    });
    
    it('name', async () => {
        const name = await contract.name().should.be.fulfilled;
        name.should.be.equal("CryptoSoccerPlayers");
    });

    it('tokenURI of unexistend player', async () => {
        await contract.tokenURI(0).should.be.rejected;
        await contract.tokenURI(1).should.be.rejected;
    });

    it('tokenURI of existent player', async () => {
        const id = 1; 
        await contract.mintWithName(accounts[0], id, "player").should.be.fulfilled;
        await contract.setTokensURI(URI)
        const uri = await contract.tokenURI(id).should.be.fulfilled;
        uri.should.be.equal(URI + "?state=0");
    });
});
