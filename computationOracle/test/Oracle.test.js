require('chai')
    .use(require('chai-as-promised'))
    .should();

const Oracle = artifacts.require('Oracle');

contract('Oracle', (accounts) => {
    it('correct deployed', async () => {
        const amount = 10;
        const oracle = await Oracle.new(amount);
        oracle.should.not.equal(null);
        const result = await oracle.stackAmount();
        result.toNumber().should.equal(amount);
    });

    it('deploy with 0 deposit required', async () => {
        await Oracle.new(0).should.be.rejected;
    });

    it('register solver with correct amount', async () => {
        const amount = 2;
        const oracle = await Oracle.new(amount);
        await oracle.registerSolver({value: amount}).should.be.fulfilled;
        const result = await oracle.solvers(accounts[0]);
        result.toNumber().should.equal(amount);
    });

    it('register solver with wrong amount', async () => {
        const amount = 2;
        const oracle = await Oracle.new(amount);
        await oracle.registerSolver({value: amount-1}).should.be.rejected;
        await oracle.registerSolver({value: amount+1}).should.be.rejected;
    });

    it('register twice the same solver', async () => {
        const amount = 2;
        const oracle = await Oracle.new(amount);
        await oracle.registerSolver({value: amount});
        await oracle.registerSolver({value: amount}).should.be.rejected;
    });

    it('unregister not registered solver', async () =>{
        const amount = 2;
        const oracle = await Oracle.new(amount);
        await oracle.unregisterSolver().should.be.rejected;
    })

    it('unregister solver', async () => {
        const amount = 2;
        const oracle = await Oracle.new(amount);
        await oracle.registerSolver({ value: amount });
        await oracle.unregisterSolver().should.be.fulfilled;
        const result = await oracle.solvers(accounts[0]);
        result.toNumber().should.equal(0);
    });
});