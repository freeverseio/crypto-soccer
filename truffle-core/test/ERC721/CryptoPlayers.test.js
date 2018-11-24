const { shouldBehaveLikeERC721 } = require('openzeppelin-solidity/test/token/ERC721/ERC721.behavior');
const BigNumber = web3.BigNumber;

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMock');

contract('CryptoPlayers', ([_, creator, ...accounts]) => {
    beforeEach(async function () {
        this.token = await CryptoPlayers.new({ from: creator });
    });

    shouldBehaveLikeERC721(creator, creator, accounts);

    it('deployment', async () => {
        await CryptoPlayers.new().should.be.fulfilled;
    });

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new();
        await contract.name().should.eventually.equal("CryptoSoccerPlayers");
        await contract.symbol().should.eventually.equal("CSP");
    });

    it('get state', async () => {
        const contract = await CryptoPlayers.new();
        const tokenId = 1;
        await contract.mint(accounts[0], tokenId).should.be.fulfilled;
        const result = await contract.getState(tokenId).should.be.fulfilled;
        result.toNumber().should.be.equal(999);
    });
});
