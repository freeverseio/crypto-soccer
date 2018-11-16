require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeams');

contract('CryptoTeams', (accounts) => {
    const name = "name";
    const symbol = "symbol";

    it('deployment', async () => {
        await CryptoTeams.new(name, symbol).should.be.fulfilled;
    });

    it('mint a team', async () => {
        const contract = await CryptoTeams.new(name, symbol).should.be.fulfilled;
        let supply = await contract.totalSupply().should.be.fulfilled;
        supply.toNumber().should.be.equal(0);
        const tokenId = 1;
        const teamName = "panzerotto";
        await contract.mint(accounts[0], tokenId, teamName).should.be.fulfilled;
        supply = await contract.totalSupply().should.be.fulfilled;
        supply.toNumber().should.be.equal(1);
    });
});
