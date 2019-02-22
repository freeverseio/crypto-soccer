require('chai')
    .use(require('chai-as-promised'))
    .should();

const Scores = artifacts.require('Scores');

contract('Scores', (accounts) => {
    let scores = null;

    beforeEach(async () => {
        scores = await Scores.new().should.be.fulfilled;
    });

    it('is valid', async () => {
        let result = await scores.isValid([5]).should.be.fulfilled;
        result.should.be.equal(false);
    });
})