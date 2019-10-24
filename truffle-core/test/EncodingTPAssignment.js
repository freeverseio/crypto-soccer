const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingTPAssignment = artifacts.require('EncodingTPAssignment');

contract('EncodingTPAssignment', (accounts) => {

    beforeEach(async () => {
        encoding = await EncodingTPAssignment.new().should.be.fulfilled;
    });
    
    it('encode and decode matchlog', async () =>  {

    });
    

});