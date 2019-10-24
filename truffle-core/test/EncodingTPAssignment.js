const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingTPAssignment = artifacts.require('EncodingTPAssignment');

contract('EncodingTPAssignment', (accounts) => {

    beforeEach(async () => {
        encoding = await EncodingTPAssignment.new().should.be.fulfilled;
        MAX_WEIGHT = await encoding.MAX_WEIGHT().should.be.fulfilled;
        MAX_WEIGHT = MAX_WEIGHT.toNumber();
        MIN_WEIGHT = await encoding.MIN_WEIGHT().should.be.fulfilled;
        MIN_WEIGHT = MIN_WEIGHT.toNumber();
    });
    it('encode fails if sum is not correct', async () =>  {
        weights =  Array.from(new Array(25), (x,i) => 3*i % 14);
        specialPlayer = 21;
        result = await encoding.encodeTP(weights, specialPlayer).should.be.rejected;
    });
    
    it('encode and decode matchlog', async () =>  {
        weights = Array.from(new Array(25), (x,i) => 3*i % 14);
        specialPlayer = 21;
        // make sure they sum to MAX_WEIGHT:
        for (bucket = 0; bucket < 5; bucket++){
            sum4 = 0;
            for (sk = 5 * bucket; sk < (5 * bucket + 4); sk++) {
                sum4 += weights[sk];
            }
            weights[5 * bucket + 4] = MAX_WEIGHT - sum4;
        }        
        
        expectedWeights = Array.from(weights, (x,i) => x + MIN_WEIGHT);
        result = await encoding.encodeTP(weights, specialPlayer).should.be.fulfilled;
        decoded = await encoding.decodeTP(result).should.be.fulfilled;
        for (bucket = 0; bucket < 5; bucket++){
            sum = 0;
            for (sk = 0; sk < 5; sk++) {
                decoded.skillWeights[5*bucket + sk].toNumber().should.be.equal(expectedWeights[5*bucket + sk]);
                sum += decoded.skillWeights[5*bucket + sk].toNumber();
            }
            (0*decoded.specialPlayer.toNumber() + sum).should.be.equal(100);
        }
        decoded.specialPlayer.toNumber().should.be.equal(specialPlayer);
    });
});