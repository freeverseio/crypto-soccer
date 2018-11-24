const { shouldBehaveLikeERC721 } = require('openzeppelin-solidity/test/token/ERC721/ERC721.behavior');
const BigNumber = web3.BigNumber;

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const TeamFactory = artifacts.require('TeamFactory');
const CryptoTeams = artifacts.require('CryptoTeamsMock');

contract('CryptoTeams', ([_, creator, ...accounts]) => {
    beforeEach(async function () {
        const teamFactory = await TeamFactory.new().should.be.fulfilled;
        this.token = await CryptoTeams.new(teamFactory.address, { from: creator });
    });

    shouldBehaveLikeERC721(creator, creator, accounts);
});
