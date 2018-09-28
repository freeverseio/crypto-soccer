require('chai')
    .use(require('chai-as-promised'))
    .should();

const Oracle = artifacts.require('Oracle');

contract('Oracle', (accounts) => {
    it('correct deployed', async () => {
        const oracle = await Oracle.new();
        oracle.should.not.equal(null);
    });

    it('registerSolver with 0 value fails', async () => {
        const oracle = await Oracle.new();
        const result = await oracle.registerSolver.call();
        result.should.equal(false);
    });
});