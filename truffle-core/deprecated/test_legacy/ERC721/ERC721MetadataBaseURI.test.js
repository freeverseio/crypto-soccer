require('chai')
    .use(require('chai-as-promised'))
    .should();

const ERC721MetadataBaseURI = artifacts.require('ERC721MetadataBaseURIMock');

contract('ERC721MetadataBaseURI', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await ERC721MetadataBaseURI.new("name", "symbol").should.be.fulfilled;
    })

   it('symbol', async () => {
        const symbol = await contract.symbol().should.be.fulfilled;
        symbol.should.be.equal("symbol");
    });
    
    it('name', async () => {
        const name = await contract.name().should.be.fulfilled;
        name.should.be.equal("name");
    });

    it('initial base token URI', async () => {
        const uri = await contract.getBaseTokenURI().should.be.fulfilled;
        uri.should.be.equal('');
    });

    it('tokenURI of unexistend token', async () => {
        await contract.tokenURI(0).should.be.rejected;
        await contract.tokenURI(1).should.be.rejected;
    });

    it('tokenURI of existent player', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        await contract.setBaseTokenURI("URI").should.be.fulfilled;
        const uri = await contract.tokenURI(id).should.be.fulfilled;
        uri.should.be.equal("URI" + id.toString(10));
    });

    it('set URI without being URIer', async () => {
        await contract.renounceURIer().should.be.fulfilled;
        await contract.setBaseTokenURI("URI").should.be.rejected;
    });
});
