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
    
    it('sorts arrays of 14 numbers', async () =>  {
        data =      [4, 7, 3, 1, 12, 9, 5, 3, 1, 6, 10, 13, 11, 11];
        expected =  [1, 1, 3, 3, 4, 5, 6, 7, 9, 10, 11, 11, 12, 13];
        result = await sort.sort14(data).should.be.fulfilled;
        for (i = 0; i < 14; i++) {
            result[i].toNumber().should.be.equal(expected[i]);
        }
    });
    
    it('sorts idxs of 8 numbers', async () =>  {
        data =          [4, 7, 3, 1, 12, 9, 5, 3];
        expectedIdxs =  [3, 2, 7, 0, 6, 1, 5, 4];
        idxs = Array.from(new Array(8), (x,i) => i);
        result = await sortIdxs.sortIdxs(data, idxs).should.be.fulfilled;
        for (i = 0; i < 8; i++) {
            // console.log(result[i].toNumber())
            result[i].toNumber().should.be.equal(expectedIdxs[i]);
        }
    });
    
});