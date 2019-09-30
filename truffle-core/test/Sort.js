const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const Sort = artifacts.require('Sort');

contract('Sort', (accounts) => {

    beforeEach(async () => {
        sort = await Sort.new().should.be.fulfilled;
    });
    
    it('encodeTactics', async () =>  {
        expected = [1, 1, 3, 3, 4, 5, 6, 7, 9, 10, 12];
        result = await sort.sort11([4,7,3,1,12,9,5,3,1,6,10]).should.be.fulfilled;
        for (i = 0; i < 11; i++) {
            result[i].toNumber().should.be.equal(expected[i]);
        }
    });
    
});