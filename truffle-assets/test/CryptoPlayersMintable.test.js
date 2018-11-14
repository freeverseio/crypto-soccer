require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayersMintable = artifacts.require('CryptoPlayersMintable');

contract('CryptoPlayersMintableMintable', (accounts) => {
    const name = "name";
    const symbol = "symbol";
    const CID = "QmUC4KA1Vi3DizRrTj9Z4uyrL6a7zjS7wNnvR5iNzYALSh";

    // it('deployment', async () => {
    //     await CryptoPlayersMintable.new(name, symbol, CID).should.be.fulfilled;
    // });

    // it('check name and symbol', async () => {
    //     const contract = await CryptoPlayersMintable.new(name, symbol, CID);
    //     await contract.name().should.eventually.equal(name);
    //     await contract.symbol().should.eventually.equal(symbol);
    // });

    // it('get state', async () => {
    //     const contract = await CryptoPlayersMintable.new(name, symbol, CID);
    //     const tokenId = 1;
    //     await contract.mint(accounts[0], tokenId).should.be.fulfilled;
    //     const result = await contract.getState(tokenId).should.be.fulfilled;
    //     result.toNumber().should.not.be.equal('0');
    // });

    // it('get URI', async () => {
    //     const contract = await CryptoPlayersMintable.new(name, symbol, CID);
    //     const tokenId = 1;
    //     await contract.mint(accounts[0], tokenId).should.be.fulfilled;
    //     const result = await contract.tokenURI(tokenId).should.be.fulfilled;
    //     result.should.be.equal(CID + '/?state=967199688875838827974656004');
    // });
});
