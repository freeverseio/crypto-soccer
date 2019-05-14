require('chai')
    .use(require('chai-as-promised'))
    .should();

const Contract = artifacts.require('Cronos');

contract('Cronos', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await Contract.new().should.be.fulfilled;        
    });

    it('wait => increare block', async () => {
        const before = await web3.eth.getBlockNumber().should.be.fulfilled;
        let after = await web3.eth.getBlockNumber().should.be.fulfilled;
        before.toString().should.be.equal(after.toString());
        await instance.wait().should.be.fulfilled;
        after = await web3.eth.getBlockNumber().should.be.fulfilled;
        before.toString().should.not.be.equal(after.toString());
    });

})