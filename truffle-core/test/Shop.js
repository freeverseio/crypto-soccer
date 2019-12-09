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
        tx = await shop.offerItem(skillsBoost, initStock = 32, matchesDuration = 8, uri =  "https://www.freeverse.io", leagueIds).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "ItemOffered", async (event) => {
            var ok = event.itemId.toNumber() === 0;
            for (l = 0; l < leagueIds.length; l++) ok = ok && (leagueIds[l] === event.leagueIds[l].toNumber());
            if (ok) return true;
            return event.itemId.toNumber().should.be.equal(-1);
        });
        result = await shop.getSkillsBoost(id = 0).should.be.fulfilled;
        result.should.be.bignumber.equal(skillsBoost);
        result = await shop.getMatchesDuration(id).should.be.fulfilled;
        result.toNumber().should.be.equal(matchesDuration);
        result = await shop.getInitStock(id).should.be.fulfilled;
        result.toNumber().should.be.equal(initStock);
        result = await shop.getUri(id).should.be.fulfilled;
        result.should.be.equal(uri);
        result = await shop.getChampionshipsHash(id).should.be.fulfilled;
        result.should.be.equal("0x39b03ec0c7168a018ef8b98732d567d4a036bc1ae1ab1a6563033f2236b362e1");
        
    });

});