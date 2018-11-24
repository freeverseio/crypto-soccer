require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMetadataMock');

contract('CryptoPlayersMetadata', (accounts) => {
    const CID = "QmUC4KA1Vi3DizRrTj9Z4uyrL6a7zjS7wNnvR5iNzYALSh";

    it('deployment', async () => {
        await CryptoPlayers.new(CID).should.be.fulfilled;
    });

    it('get URI', async () => {
        const contract = await CryptoPlayers.new(CID);
        const tokenId = 1;
        await contract.mint(accounts[0], tokenId).should.be.fulfilled;
        const result = await contract.tokenURI(tokenId).should.be.fulfilled;
        result.should.be.equal(CID + '/?state=967199688875838827974656004');
    });
});
