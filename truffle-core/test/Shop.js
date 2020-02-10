const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');

const Shop = artifacts.require('Shop');
const EncodingTactics = artifacts.require('EncodingTactics');

contract('Shop', (accounts) => {
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];

    const it2 = async(text, f) => {};

    beforeEach(async () => {
        shop = await Shop.new().should.be.fulfilled;
        await shop.init().should.be.fulfilled;
        encTactics = await EncodingTactics.new().should.be.fulfilled;
    });

    
    it2('encode - decode boosts', async () => {
        boosts = [62,60,19,1,23,2];
        encoded = await shop.encodeBoosts(boosts).should.be.fulfilled;
        decoded = await shop.decodeBoosts(encoded).should.be.fulfilled;
        debug.compareArrays(decoded, boosts, toNum = true, verbose = false, isBigNumber = false);
    });
    
    it2('offer item', async () => {
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

        encodedBoost = await shop.encodeBoosts(boosts).should.be.fulfilled;
        
        truffleAssert.eventEmitted(tx, "ItemOffered", (event) => {
            return event.itemId.toNumber() === 1 && 
                event.encodedBoost.toNumber() == encodedBoost &&
                event.countriesRoot.toNumber() === countriesRoot &&
                event.championshipsRoot.toNumber() === championshipsRoot &&
                event.teamsRoot.toNumber() === teamsRoot &&
                event.itemsRemaining.toNumber() === itemsRemaining &&
                event.matchesDuration.toNumber() === matchesDuration &&
                event.onlyTopInChampioniship.toNumber() === onlyTopInChampioniship &&
                event.uri === uri;
        }, "correct");
        
        await shop.reduceItemsRemaining(itemId = 1, itemsRemaining - 3).should.be.fulfilled;
        result = await shop.getItemsRemaining(itemId).should.be.fulfilled;
        result.toNumber().should.be.equal(itemsRemaining - 3);
    });
    
    it('add items to tactics', async () => {
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
        encodedBoost = await shop.encodeBoosts(boosts).should.be.fulfilled;
        
        const lineupConsecutive = Array.from(new Array(14), (x,i) => i); 
        const extraAttackNull =  Array.from(new Array(10), (x,i) => 0);
        tactics = await encTactics.encodeTactics(substitutions = [3,4,5], subsRounds = [6,7,8], lineupConsecutive, 
            extraAttackNull, tacticId442 = 2).should.be.fulfilled;
        
        tactics2 = await shop.addItemsToTactics(tactics, itemId = 1, staminaRecovery = 2).should.be.fulfilled;
        const {0: stamina, 1: id, 2: boost} = await shop.getItemsData(tactics2).should.be.fulfilled;
        stamina.toNumber().should.be.equal(staminaRecovery);
        id.toNumber().should.be.equal(itemId);
        boost.should.be.bignumber.equal(encodedBoost);
        
        
    });

});