const { shouldBehaveLikeERC721 } = require('openzeppelin-solidity/test/token/ERC721/ERC721.behavior');

const BigNumber = web3.BigNumber;
const ERC721Mock = artifacts.require('ERC721Mock.sol');

require('chai')
  .use(require('chai-bignumber')(BigNumber))
  .should();

contract('ERC721', function ([_, creator, ...accounts]) {
  beforeEach(async function () {
    this.token = await ERC721Mock.new({ from: creator });
  });

  shouldBehaveLikeERC721(creator, creator, accounts);
});