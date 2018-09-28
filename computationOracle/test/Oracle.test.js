require('chai')
    .use(require('chai-as-promised'))
    .should();

const Oracle = artifacts.require('Oracle');

contract('Oracle', (accounts) => {
    it('correct deployed', async () => {
        const oracle = await Oracle.new(0);
        oracle.should.not.equal(null);
    });

    it('registerSolver with correct amount', async () => {
        const amount = 100;
        const oracle = await Oracle.new(amount);
        await oracle.registerSolver.call({value: amount}).should.be.fulfilled;
    });

    it('registerSolver with wrong amount', async () => {
        const amount = 100;
        const oracle = await Oracle.new(amount);
        await oracle.registerSolver.call({value: amount-1}).should.be.rejected;
        await oracle.registerSolver.call({value: amount+1}).should.be.rejected;
    });
});