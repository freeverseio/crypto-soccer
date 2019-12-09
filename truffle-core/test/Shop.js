const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');

const Shop = artifacts.require('Shop');

contract('Shop', (accounts) => {
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];

    const it2 = async(text, f) => {};

    beforeEach(async () => {
        shop = await Shop.new().should.be.fulfilled;
    });

    it('offer item', async () => {
        skillsBoost = await shop.encodeSkillsBoost([5,7,8,11,23]).should.be.fulfilled;
        leagueIds = Array.from(new Array(32), (x,i) => 312312321+i);
        tx = await shop.offerItem(skillsBoost, stock = 32, matchesDuration = 8, "https://www.freeverse.io", leagueIds).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "ItemOffered", async (event) => {
            return event.itemId.toNumber().should.be.equal(0);
            // return true;
            // return event.playerId.should.be.bignumber.equal(playerId) && currentTeam.should.be.bignumber.equal(teamId2);
        });

    });

});