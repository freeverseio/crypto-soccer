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
        await shop.init().should.be.fulfilled;
    });

    
    it('encode - decode boosts', async () => {
        boosts = [62,60,19,1,23,2];
        encoded = await shop.encodeBoosts(boosts).should.be.fulfilled;
        decoded = await shop.decodeBoosts(encoded).should.be.fulfilled;
        debug.compareArrays(decoded, boosts, toNum = true, verbose = false, isBigNumber = false);
    });
    
    it('offer item', async () => {
        tx = await shop.offerItem(
            boosts = [62,60,19,1,23,1],
            countriesRoot = 0,
            championshipsRoot = 0,
            teamsRoot = 0,
            itemsRemaining = 5432,
            matchesDuration = 7,
            onlyTopInChampioniship = 3,
            uri =  "https://www.freeverse.io"
        ).should.be.rejected;

        tx = await shop.offerItem(
            boosts = [32,30,19,1,23,1],
            countriesRoot = 0,
            championshipsRoot = 0,
            teamsRoot = 0,
            itemsRemaining = 5432,
            matchesDuration = 7,
            onlyTopInChampioniship = 3,
            uri =  "https://www.freeverse.io"
        ).should.be.fulfilled;

        truffleAssert.eventEmitted(tx, "ItemOffered", (event) => {
            return event.itemId.toNumber() === 1 && 
                event.countriesRoot.toNumber() === countriesRoot &&
                event.championshipsRoot.toNumber() === championshipsRoot &&
                event.teamsRoot.toNumber() === teamsRoot &&
                event.itemsRemaining.toNumber() === itemsRemaining &&
                event.matchesDuration.toNumber() === matchesDuration &&
                event.onlyTopInChampioniship.toNumber() === onlyTopInChampioniship &&
                event.uri === uri;
        }, "correct");
        // result = await shop.getSkillsBoost(id = 0).should.be.fulfilled;
        // result.should.be.bignumber.equal(skillsBoost);
        // result = await shop.getMatchesDuration(id).should.be.fulfilled;
        // result.toNumber().should.be.equal(matchesDuration);
        // result = await shop.getInitStock(id).should.be.fulfilled;
        // result.toNumber().should.be.equal(initStock);
        // result = await shop.getUri(id).should.be.fulfilled;
        // result.should.be.equal(uri);
        // result = await shop.getChampionshipsHash(id).should.be.fulfilled;
        // result.should.be.equal("0x39b03ec0c7168a018ef8b98732d567d4a036bc1ae1ab1a6563033f2236b362e1");
        
    });

});