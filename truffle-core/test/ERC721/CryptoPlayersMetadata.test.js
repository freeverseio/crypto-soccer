require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMetadataMock');

contract('CryptoPlayersMetadata', (accounts) => {
    const URI = "QmUC4KA1Vi3DizRrTj9Z4uyrL6a7zjS7wNnvR5iNzYALSh";

    it('deployment', async () => {
        await CryptoPlayers.new().should.be.fulfilled;
    });

    it('set base URI', async () => {
        const contract = await CryptoPlayers.new();
        await contract.setBaseURI(URI).should.be.fulfilled;
        const result = await contract.getBaseURI().should.be.fulfilled;
        result.should.be.equal(URI);
    });

    it('get URI', async () => {
        const contract = await CryptoPlayers.new();
        await contract.setBaseURI(URI).should.be.fulfilled;
        const tokenId = 1;
        await contract.mint(accounts[0], tokenId).should.be.fulfilled;
        const result = await contract.tokenURI(tokenId).should.be.fulfilled;
        result.should.be.equal(URI + '/?state=999');
    });
});
