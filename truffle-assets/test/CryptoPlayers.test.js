require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    const name = "name";
    const symbol = "symbol";

    it('deployment', async () => {
        await CryptoPlayers.new(name, symbol).should.be.fulfilled;
    });

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new(name, symbol);
        await contract.name().should.eventually.equal(name);
        await contract.symbol().should.eventually.equal(symbol);
    });

    it('mint a player', async () => {
        const contract = await CryptoPlayers.new(name, symbol);
        const tokenId = 1;
        const tokenURI = "http://russo";
        await contract.mintWithTokenURI(accounts[0], tokenId, tokenURI).should.be.fulfilled;
        await contract.tokenURI(tokenId).should.eventually.equal(tokenURI);
    });
});
