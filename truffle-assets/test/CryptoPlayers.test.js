require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMock');

contract('CryptoPlayers', (accounts) => {
    const name = "name";
    const symbol = "symbol";
    const CID = "http://freeverse.io";

    it('deployment', async () => {
        await CryptoPlayers.new(name, symbol, CID).should.be.fulfilled;
    });

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new(name, symbol, CID);
        await contract.name().should.eventually.equal(name);
        await contract.symbol().should.eventually.equal(symbol);
    });

    it('get URI', async () => {
        const contract = await CryptoPlayers.new(name, symbol, CID);
        const tokenId = 1;
        const state = 0;
        await contract.mint(accounts[0], tokenId, state).should.be.fulfilled;
        const URI = await contract.tokenURI(tokenId).should.be.fulfilled;
        URI.should.be.equal(CID);
    });

    // it('mint a player', async () => {
    //     const contract = await CryptoPlayers.new(name, symbol);
    //     const tokenId = 1;
    //     const tokenURI = "http://russo";
    //     await contract.mintWithTokenURI(accounts[0], tokenId, tokenURI).should.be.fulfilled;
    //     await contract.tokenURI(tokenId).should.eventually.equal(tokenURI);
    // });
});
