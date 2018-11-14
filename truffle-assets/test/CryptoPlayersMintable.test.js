require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayersMintable = artifacts.require('CryptoPlayersMintable');

contract('CryptoPlayersMintableMintable', (accounts) => {
    it('deployment', async () => {
        await CryptoPlayersMintable.new().should.be.fulfilled;
    });

    it('mint', async () => {
        const contract = await CryptoPlayersMintable.new();
        const tokenId = 1;
        await contract.mint(accounts[0], tokenId).should.be.fulfilled;
        const result = await contract.getState(tokenId).should.be.fulfilled;
        result.toNumber().should.be.equal(9.671996888758388e+26);
    });
});
