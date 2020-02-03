const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingTPAssignment = artifacts.require('EncodingTPAssignment');

contract('EncodingTPAssignment', (accounts) => {

    beforeEach(async () => {
        encoding = await EncodingTPAssignment.new().should.be.fulfilled;
        MAX_PERCENT = await encoding.MAX_PERCENT().should.be.fulfilled;
        MAX_PERCENT = MAX_PERCENT.toNumber();
        MIN_PERCENT = await encoding.MIN_PERCENT().should.be.fulfilled;
        MIN_PERCENT = MIN_PERCENT.toNumber();
    });
    
    it('encode fails if sum is not correct', async () =>  {
        specialPlayer = 21;
        TP = 40;
        TPperSkill =  Array.from(new Array(25), (x,i) => Math.floor(TP/5));
        result = await encoding.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        // value too small:
        TPperSkill[2] = 1;
        result = await encoding.encodeTP(TP, TPperSkill, specialPlayer).should.be.rejected;
        // sum too large:
        TPperSkill =  Array.from(new Array(25), (x,i) => 1 + Math.floor(TP/5));
        result = await encoding.encodeTP(TP, TPperSkill, specialPlayer).should.be.rejected;
        // special player can have extra 10 percent:
        TPperSkill =  Array.from(new Array(25), (x,i) => Math.floor(TP/5));
        // special player belongs to idx: 20, 21, 22, 23, 24
        // having an extra 10% he should be able to reach 44 points => (8,8,8,8,8) -> 9,8,9,9,9
        TPperSkill[20] = 9;
        TPperSkill[21] = 8;
        TPperSkill[22] = 9;
        TPperSkill[23] = 9;
        TPperSkill[24] = 9;
        result = await encoding.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        TPperSkill[24] = 10;
        result = await encoding.encodeTP(TP, TPperSkill, specialPlayer).should.be.rejected;
    });

    it('encode and decode TP assignment', async () =>  {
        specialPlayer = 21;
        TP = 40;
        TPperSkill = Array.from(new Array(25), (x,i) => 3 + 3*i % 5);
        // make sure they sum to TP:
        for (bucket = 0; bucket < 5; bucket++){
            sum4 = 0;
            for (sk = 5 * bucket; sk < (5 * bucket + 4); sk++) {
                sum4 += TPperSkill[sk];
            }
            TPperSkill[5 * bucket + 4] = TP - sum4;
        }        
        result = await encoding.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        decoded = await encoding.decodeTP(result).should.be.fulfilled;
        for (bucket = 0; bucket < 5; bucket++){
            sum = 0;
            for (sk = 0; sk < 5; sk++) {
                decoded.TPperSkill[5*bucket + sk].toNumber().should.be.equal(TPperSkill[5*bucket + sk]);
                sum += decoded.TPperSkill[5*bucket + sk].toNumber();
            }
            (0*decoded.specialPlayer.toNumber() + sum).should.be.equal(TP);
        }
        decoded.specialPlayer.toNumber().should.be.equal(specialPlayer);
        decoded.TP.toNumber().should.be.equal(TP);
    });
});