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
        await instance.encodeScore(0xff, 2).should.be.fulfilled;
        await instance.encodeScore(2, 0xff).should.be.fulfilled;
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

    it('fill a day scores', async () => {
        await instance.addToDayScores([0xffff], 0x0101).should.be.rejected;
        await instance.addToDayScores([], 0xffff).should.be.rejected;
        let scores = [];
        let score = await instance.encodeScore(3, 0).should.be.fulfilled;
        scores = await instance.addToDayScores(scores, score).should.be.fulfilled;
        score = await instance.encodeScore(1, 2).should.be.fulfilled;
        scores = await instance.addToDayScores(scores, score).should.be.fulfilled;
        score = await instance.encodeScore(0, 0).should.be.fulfilled;
        scores = await instance.addToDayScores(scores, score).should.be.fulfilled;
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
    });

    it('fill tournament scores', async () => {
        let scores = [];
        scores = await instance.addToTournamentScores(scores, [0x0201, 0x0101]).should.be.fulfilled;
        scores = await instance.addToTournamentScores(scores, [0x0101, 0x0001]).should.be.fulfilled;
        scores = await instance.addToTournamentScores(scores, [0x0001, 0x0004]).should.be.fulfilled;
        await instance.addToTournamentScores(scores, [0xffff, 0x0101]).should.be.rejected;
        scores = await instance.addToTournamentScores(scores, []).should.be.fulfilled;
        scores.length.should.be.equal(9);
        scores[0].toNumber().should.be.equal(0x0201);
        scores[1].toNumber().should.be.equal(0x0101);
        scores[2].toNumber().should.be.equal(0xffff);
        scores[3].toNumber().should.be.equal(0x0101);
        scores[4].toNumber().should.be.equal(0x0001);
        scores[5].toNumber().should.be.equal(0xffff);
        scores[6].toNumber().should.be.equal(0x0001);
        scores[7].toNumber().should.be.equal(0x0004);
        scores[8].toNumber().should.be.equal(0xffff);
    });

    it('count days', async () => {
        let result = await instance.countDaysInTournamentScores([]).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await instance.countDaysInTournamentScores([2, 3]).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        let leagueState = await instance.addToTournamentScores([2, 3], [4]).should.be.fulfilled;
        result = await instance.countDaysInTournamentScores(leagueState).should.be.fulfilled;
        result.toNumber().should.be.equal(2);
        leagueState = await instance.addToTournamentScores(leagueState, [4, 5, 6, 6]).should.be.fulfilled;
        result = await instance.countDaysInTournamentScores(leagueState).should.be.fulfilled;
        result.toNumber().should.be.equal(3);
    });
})