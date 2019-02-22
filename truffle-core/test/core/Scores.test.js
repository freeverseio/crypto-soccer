require('chai')
    .use(require('chai-as-promised'))
    .should();

const Scores = artifacts.require('Scores');

contract('Scores', (accounts) => {
    let scores = null;
    let divider = null;

    beforeEach(async () => {
        scores = await Scores.new().should.be.fulfilled;
        divider = await scores.DIVIDER().should.be.fulfilled;
    });

    it('is valid', async () => {
        let result = await scores.isValid([5]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await scores.isValid([]).should.be.fulfilled;
        result.should.be.equal(true);
        result = await scores.isValid([divider, 4]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await scores.isValid([4, divider]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await scores.isValid([4, 2, divider, divider, 3, 4]).should.be.fulfilled;
        result.should.be.equal(false);
        result = await scores.isValid([4, 2, divider, 4, 3, 3, 2, divider, 3, 4]).should.be.fulfilled;
        result.should.be.equal(true);
    });
})