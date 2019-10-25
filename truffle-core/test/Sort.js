const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const SortValues = artifacts.require('SortValues');
const SortIdxs = artifacts.require('SortIdxs');

contract('SortValues', (accounts) => {

    const it2 = async(text, f) => {};

    beforeEach(async () => {
        sort = await SortValues.new().should.be.fulfilled;
        sortIdxs = await SortIdxs.new().should.be.fulfilled;
    });
    
    it2('sorts arrays of 14 numbers', async () =>  {
        data =      [4, 7, 3, 1, 12, 9, 5, 3, 1, 6, 10, 13, 11, 11];
        expected =  [1, 1, 3, 3, 4, 5, 6, 7, 9, 10, 11, 11, 12, 13];
        result = await sort.sort14(data).should.be.fulfilled;
        for (i = 0; i < 14; i++) {
            result[i].toNumber().should.be.equal(expected[i]);
        }
    });
    
    it('sorts idxs of 14 numbers', async () =>  {
        data =          [4, 7, 3, 1, 12, 9, 5, 3, 1, 6, 10, 13, 11, 11];
        expectedVals =  [1, 1, 3, 3, 4, 5, 6, 7, 9, 10, 11, 11, 12, 13];
        expectedIdxs =  [ 3, 8, 7, 2, 0, 6, 9, 1, 5, 10, 12, 13, 4, 11];
        result = await sortIdxs.sort14(data).should.be.fulfilled;
        ex = [];
        for (i = 0; i < 14; i++) {
            result[i].toNumber().should.be.equal(expectedIdxs[i]);
        }
    });
    
});