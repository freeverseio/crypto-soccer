const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const SortValues = artifacts.require('SortValues');
const SortIdxs = artifacts.require('SortIdxs');

contract('SortValues', (accounts) => {

    const it2 = async(text, f) => {};

    function compareArrays(result, expected, toNum = true, verbose = false) {
        verb = [];
        for (i = 0; i < expected.length; i++) {
            if (toNum) verb.push(result[i].toNumber());
            else verb.push(result[i]);
            if (!verbose) {
                if (toNum) result[i].toNumber().should.be.equal(expected[i]);
                else result[i].should.be.equal(expected[i]);
            }            
        }
        if (verbose) console.log(verb);
    }
    
    beforeEach(async () => {
        sort = await SortValues.new().should.be.fulfilled;
        sortIdxs = await SortIdxs.new().should.be.fulfilled;
    });
    
    it2('sorts arrays of 14 numbers', async () =>  {
        data =      [4, 7, 3, 1, 12, 9, 5, 3, 1, 6, 10, 13, 11, 11];
        expected =  [13, 12, 11, 11, 10, 9, 7, 6, 5, 4, 3, 3, 1, 1];
        result = await sort.sort14(data).should.be.fulfilled;
        compareArrays(result, expected, toNum = true, verbose = false);
    });
    
    it('sorts idxs of 8 numbers', async () =>  {
        data =          [4, 7, 3, 1, 12, 9, 5, 3];
        expectedIdxs =  [ 4, 5, 1, 6, 0, 2, 7, 3 ];
        idxs = Array.from(new Array(8), (x,i) => i);
        result = await sortIdxs.sortIdxs(data, idxs).should.be.fulfilled;
        compareArrays(result, expectedIdxs, toNum = true, verbose = false);
    });
    
});