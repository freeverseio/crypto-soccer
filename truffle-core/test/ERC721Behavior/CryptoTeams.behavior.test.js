const { shouldBehaveLikeERC721 } = require('openzeppelin-solidity/test/token/ERC721/ERC721.behavior');
const BigNumber = web3.BigNumber;

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const CryptoTeams = artifacts.require('CryptoTeamsMock');

contract('CryptoTeams', ([_, creator, ...accounts]) => {
    beforeEach(async function () {
        this.token = await CryptoTeams.new({ from: creator });
    });

    shouldBehaveLikeERC721(creator, creator, accounts);
});
