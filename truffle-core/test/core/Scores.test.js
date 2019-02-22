require('chai')
    .use(require('chai-as-promised'))
    .should();

const Scores = artifacts.require('Scores');

contract('Scores', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await Scores.new().should.be.fulfilled;
    });

    it('encode scores', async () => {
        await instance.encodeScore(0xff, 2).should.be.rejected;
        await instance.encodeScore(2, 0xff).should.be.rejected;
        await instance.encodeScore(0xff, 0xff).should.be.rejected;
        const score = await instance.encodeScore(0x01,0x02).should.be.fulfilled;
        score.toNumber().should.be.equal(0x0102);
    });

    it('decode', async () => {
        await instance.decodeScore(0xffff).should.be.rejected;
        const result = await instance.decodeScore(0x0102).should.be.fulfilled;
        result.home.toNumber().should.be.equal(0x01);
        result.visitor.toNumber().should.be.equal(0x02);
    })

    // it('is valid', async () => {
    //     let result = await scores.isValid([5]).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     result = await scores.isValid([]).should.be.fulfilled;
    //     result.should.be.equal(true);
    //     result = await scores.isValid([divider, 4]).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     result = await scores.isValid([4, divider]).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     result = await scores.isValid([4, 2, divider, divider, 3, 4]).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     result = await scores.isValid([4, 2, divider, 4, 3, 3, 2, 3, 4]).should.be.fulfilled;
    //     result.should.be.equal(true);
    //     result = await scores.isValid([4, 2, divider, 4, 3, 3, 2, divider, 3, 4]).should.be.fulfilled;
    //     result.should.be.equal(true);
    // });

    // it('count days', async () => {
    //     let result = await scores.scoresCountDays([]).should.be.fulfilled;
    //     result.toNumber().should.be.equal(0);
    //     result = await scores.scoresCountDays([2, 3]).should.be.fulfilled;
    //     result.toNumber().should.be.equal(1);
    //     result = await scores.scoresCountDays([2, 3, divider, 4 ,5]).should.be.fulfilled;
    //     result.toNumber().should.be.equal(1);
    //     result = await scores.scoresCountDays([2, 3]).should.be.fulfilled;
    //     result.toNumber().should.be.equal(1);
    // })
})