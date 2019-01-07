require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMetadata');

contract('CryptoPlayersMetadata', (accounts) => {
    let contract = null;

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

    it('initial base token URI', async () => {
        const uri = await contract.getBaseTokenURI().should.be.fulfilled;
        uri.should.be.equal('');
    });

    it('tokenURI of unexistend player', async () => {
        await contract.tokenURI(0).should.be.rejected;
        await contract.tokenURI(1).should.be.rejected;
    });

    it('tokenURI of existent player', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setBaseTokenURI("URI").should.be.fulfilled;
        const uri = await contract.tokenURI(id).should.be.fulfilled;
        const genome = await contract.getGenome(id).should.be.fulfilled;
        uri.should.be.equal("URI/" + id.toString(10));
    });

    it('set URI without being URIer', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        await contract.renounceURIer().should.be.fulfilled;
        await contract.setBaseTokenURI("URI").should.be.rejected;
    });
});
