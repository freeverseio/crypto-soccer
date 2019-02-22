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
        const score = await instance.decodeScore(0x0102).should.be.fulfilled;
        score.home.toNumber().should.be.equal(0x01);
        score.visitor.toNumber().should.be.equal(0x02);
    });

    it('fill a day match', async () => {
        await instance.addScore([0xffff], 0x0101).should.be.rejected;
        await instance.addScore([], 0xffff).should.be.rejected;
        let scores = [];
        let score = await instance.encodeScore(3, 0).should.be.fulfilled;
        scores = await instance.addScore(scores, score).should.be.fulfilled;
        score = await instance.encodeScore(1, 2).should.be.fulfilled;
        scores = await instance.addScore(scores, score).should.be.fulfilled;
        score = await instance.encodeScore(0, 0).should.be.fulfilled;
        scores = await instance.addScore(scores, score).should.be.fulfilled;
        scores.length.should.be.equal(3);
        score = await instance.decodeScore(scores[0]).should.be.fulfilled;
        score.home.toNumber().should.be.equal(3);
        score.visitor.toNumber().should.be.equal(0);
        score = await instance.decodeScore(scores[1]).should.be.fulfilled;
        score.home.toNumber().should.be.equal(1);
        score.visitor.toNumber().should.be.equal(2);
        score = await instance.decodeScore(scores[2]).should.be.fulfilled;
        score.home.toNumber().should.be.equal(0);
        score.visitor.toNumber().should.be.equal(0);
    });

    it('check day match validity', async () => {
        let result = await instance.isValidDayScores([]);
        result.should.be.equal(true);
        result = await instance.isValidDayScores([0x0201, 0x0101, 0x0000]);
        result.should.be.equal(true);
        result = await instance.isValidDayScores([0x0201, 0xffff]);
        result.should.be.equal(false);
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