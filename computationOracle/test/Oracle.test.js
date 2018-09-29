require('chai')
    .use(require('chai-as-promised'))
    .should();

const Oracle = artifacts.require('Oracle');

contract('Oracle', (accounts) => {
    it('correct deployed', async () => {
        const deposit = 10;
        const oracle = await Oracle.new(deposit);
        oracle.should.not.equal(null);
        const result = await oracle.deposit();
        result.toNumber().should.equal(deposit);
    });

    it('deploy with 0 deposit required', async () => {
        await Oracle.new(0).should.be.rejected;
    });

    it('register solver with correct deposit', async () => {
        const deposit = 2;
        const oracle = await Oracle.new(deposit);
        await oracle.registerSolver({value: deposit}).should.be.fulfilled;
        const result = await oracle.solvers(accounts[0]);
        result.toNumber().should.equal(deposit);
    });

    it('register solver with wrong deposit', async () => {
        const deposit = 2;
        const oracle = await Oracle.new(deposit);
        await oracle.registerSolver({value: deposit-1}).should.be.rejected;
        await oracle.registerSolver({value: deposit+1}).should.be.rejected;
    });

    it('register twice the same solver', async () => {
        const deposit = 2;
        const oracle = await Oracle.new(deposit);
        await oracle.registerSolver({value: deposit});
        await oracle.registerSolver({value: deposit}).should.be.rejected;
    });

    it('unregister not registered solver', async () =>{
        const deposit = 2;
        const oracle = await Oracle.new(deposit);
        await oracle.unregisterSolver().should.be.rejected;
    })

    it('unregister solver', async () => {
        const deposit = 2;
        const oracle = await Oracle.new(deposit);
        await oracle.registerSolver({ value: deposit });
        await oracle.unregisterSolver().should.be.fulfilled;
        const result = await oracle.solvers(accounts[0]);
        result.toNumber().should.equal(0);
    });
});