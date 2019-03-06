require('chai')
    .use(require('chai-as-promised'))
    .should();

const Scores = artifacts.require('LeaguesScore');

contract('LeaguesScore', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await Scores.new().should.be.fulfilled;
    });

    it('encode scores', async () => {
        let score = await instance.encodeScore(0xff, 2).should.be.fulfilled;
        score.toNumber().should.be.equal(0xff02);
        score = await instance.encodeScore(0x01,0x02).should.be.fulfilled;
        score.toNumber().should.be.equal(0x0102);
    });

    it('decode', async () => {
        const score = await instance.decodeScore(0x0102).should.be.fulfilled;
        score.home.toNumber().should.be.equal(0x01);
        score.visitor.toNumber().should.be.equal(0x02);
    });

    it('fill a day scores', async () => {
        let scores = await instance.scoresCreate().should.be.fulfilled;
        let score = await instance.encodeScore(3, 0).should.be.fulfilled;
        scores = await instance.scoresAppend(scores, score).should.be.fulfilled;
        score = await instance.encodeScore(1, 2).should.be.fulfilled;
        scores = await instance.scoresAppend(scores, score).should.be.fulfilled;
        score = await instance.encodeScore(0, 0).should.be.fulfilled;
        scores = await instance.scoresAppend(scores, score).should.be.fulfilled;
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
})