const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let engine = null;
    let encoding = null;
    let teamStateAll50 = null;
    const seed = 610106;
    const tactic0 = 0; // 442
    const tactic1 = 1; // 541

    const createTeamStateFromSinglePlayer = async (skills, encoding) => {
        const playerStateTemp = await encoding.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1
        ).should.be.fulfilled;

        teamState = []
        for (player = 0; player < 11; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };
    

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        encoding = engine;
        await engine.init().should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], encoding);
        });

    // it('teams get tired', async () => {
    //     const result = await engine.teamsGetTired([10,20,30,40,100], [20,40,60,80,50]).should.be.fulfilled;
    //     result[0][0].toNumber().should.be.equal(10);
    //     result[0][1].toNumber().should.be.equal(20);
    //     result[0][2].toNumber().should.be.equal(30);
    //     result[0][3].toNumber().should.be.equal(40);
    //     result[0][4].toNumber().should.be.equal(100);
    //     result[1][0].toNumber().should.be.equal(10);
    //     result[1][1].toNumber().should.be.equal(20);
    //     result[1][2].toNumber().should.be.equal(30);
    //     result[1][3].toNumber().should.be.equal(40);
    //     result[1][4].toNumber().should.be.equal(50);
    // });

    // it('play a match', async () => {
    //     let teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], encoding);
    //     const result = await engine.playMatch(seed, teamStateAll50, teamStateAll1, tactic0, tactic1).should.be.fulfilled;
    //     result[0].toNumber().should.be.equal(17);
    //     result[1].toNumber().should.be.equal(0);
    // });

    it('manages to score', async () => {
        // interface: 
        // managesToScore(uint8 nAttackers, uint[] attackersSpeed, uint[], attackersShoot, blockShoot, rndNum1,rndNum2)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        let attackersSpeed = [10,1,1];
        let attackersShoot = [10,1,1];
        let blockShoot     = 1;
        nAttackers         = attackersShoot.length;
        let result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
        blockShoot     = 1000;
        result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        // even with a super-goalkeeper, there are chances of scoring (e.g. if the rnd is super small, in this case)
        kMaxRndNumHalf = 1;
        result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('throws dice array', async () => {
        // interface: throwDiceArray(uint[] memory weights, uint rndNum)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        let result = await engine.throwDiceArray([1000,1,1],kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDiceArray([1,1000,1],kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await engine.throwDiceArray([1,1,1000],kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(2);
    });

    it('manages to shoot', async () => {
        // interface: managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        let globSkills = [[100,100,100,100,100], [1,1,1,1,1]];
        let result = await engine.managesToShoot(0,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
        result = await engine.managesToShoot(1,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        globSkills = [[1,1,1,1,1], [100,100,100,100,100]];
        result = await engine.managesToShoot(0,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        result = await engine.managesToShoot(1,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
    });


});