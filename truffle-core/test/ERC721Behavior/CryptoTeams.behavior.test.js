const { shouldBehaveLikeERC721 } = require('openzeppelin-solidity/test/token/ERC721/ERC721.behavior');
const BigNumber = web3.BigNumber;

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const Teams = artifacts.require('TeamsMock');

contract('Teams', ([_, creator, ...accounts]) => {
    const PlayersAddress = 0;

    beforeEach(async function () {
        this.token = await Teams.new(PlayersAddress, { from: creator });
    });

    shouldBehaveLikeERC721(creator, creator, accounts);
});
