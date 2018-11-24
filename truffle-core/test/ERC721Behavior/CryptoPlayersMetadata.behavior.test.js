const { shouldBehaveLikeERC721 } = require('openzeppelin-solidity/test/token/ERC721/ERC721.behavior');
const BigNumber = web3.BigNumber;

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMetadataMock');

contract('CryptoPlayersMetadata', ([_, creator, ...accounts]) => {
    beforeEach(async function () {
        this.token = await CryptoPlayers.new({ from: creator });
    });

    shouldBehaveLikeERC721(creator, creator, accounts);
});
