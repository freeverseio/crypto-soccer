using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Numerics;
using System;


namespace NewMSTestProject
{

    [TestClass]
    public class UnitTest1
    {
        [TestMethod]
        public void GetCurrentShirtNum() {  
            Serialization serial = new Serialization();
            uint shirt = 13;
            BigInteger state = new BigInteger(shirt * Math.Pow(2,43));
            uint res = serial.getCurrentShirtNum(state);  
            Assert.AreEqual(res, (uint) shirt);  
        }  

        [TestMethod]
        public void DecodePlayerState() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingStateTestData.json");
            foreach(dynamic test in tests) {
                tu.AssertDecodePlayerStateOK(test);
            }
        }  

        [TestMethod]
        public void DecodePlayerStateFromTheField() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingStateTestDataFromTheField.json");
            foreach(dynamic test in tests) {
                tu.AssertDecodePlayerStateOK(test);
            }
        }  

        [TestMethod]
        public void DecodeTeamAndPlayerIDs() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingPlayerIDsDataFromTheField.json");
            foreach(dynamic test in tests) {
                BigInteger encoded;
                bool succeeded = BigInteger.TryParse((string) test.encodedId, out encoded);
                Assert.AreEqual(true, succeeded);
                Assert.AreEqual((uint) test.timezone, serial.getTimezone(encoded));
                Assert.AreEqual((uint) test.country, serial.getCountryIdxInTZ(encoded));
                Assert.AreEqual((uint) test.val, serial.getValInCountry(encoded));
            }  
        }

        [TestMethod]
        public void DecodeTPs() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingTPsData.json");
            foreach(dynamic test in tests) {
                BigInteger encoded;
                bool succeeded = BigInteger.TryParse((string) test.encodedTPAssignment, out encoded);
                Assert.AreEqual(true, succeeded);
                (uint[] TPperSkill, uint specialPlayer, uint TP, uint err) = serial.decodeTP(encoded);

                for (int i = 0; i < 25; i++) { Assert.AreEqual((uint) test.TPperSkill[i], TPperSkill[i]); }
                Assert.AreEqual((uint) test.specialPlayer, specialPlayer);
                Assert.AreEqual((uint) test.TP, TP);
                Assert.AreEqual((uint) test.err, err);
            }
        }  

        [TestMethod]
        public void EncodeTPs() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingTPsData.json");
            foreach(dynamic test in tests) {
                BigInteger encodedExpected;
                bool succeeded = BigInteger.TryParse((string) test.encodedTPAssignment, out encodedExpected);
                Assert.AreEqual(true, succeeded); 
                uint[] TPperSkill = tu.DynamicToUintArray(test.TPperSkill);
                (BigInteger encoded, string err) = serial.encodeTP(
                    (uint) test.TP,
                    TPperSkill,
                    (uint) test.specialPlayer
                );
                Assert.AreEqual(err, "");  
                Assert.AreEqual(encoded, encodedExpected);  
            }
        }  

        [TestMethod]
        public void DecodePlayerSkills() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingSkillsTestData.json");
            foreach(dynamic test in tests) {
                BigInteger encoded;
                bool succeeded = BigInteger.TryParse((string) test.encodedSkills, out encoded);
                Assert.AreEqual(true, succeeded);
                for (int sk = 0; sk < 5; sk++) { Assert.AreEqual((uint) test.skills[sk], serial.getSkill(encoded, sk)); }
                Assert.AreEqual((uint) test.birthday, serial.getBirthDay(encoded));  
                Assert.AreEqual((bool) test.isSpecial, serial.getIsSpecial(encoded));  
                Assert.AreEqual((uint) test.potential, serial.getPotential(encoded));  
                Assert.AreEqual((uint) test.forwardness, serial.getForwardness(encoded));  
                Assert.AreEqual((uint) test.leftishness, serial.getLeftishness(encoded));  
                Assert.AreEqual((uint) test.aggressiveness, serial.getAggressiveness(encoded));  
                Assert.AreEqual((bool) test.alignedEndOfFirstHalf, serial.getAlignedEndOfFirstHalf(encoded));  
                Assert.AreEqual((bool) test.redCardLastGame, serial.getRedCardLastGame(encoded));  
                Assert.AreEqual((uint) test.gamesNonStopping, serial.getGamesNonStopping(encoded));  
                Assert.AreEqual((uint) test.injuryWeeksLeft, serial.getInjuryWeeksLeft(encoded));  
                Assert.AreEqual((bool) test.substitutedFirstHalf, serial.getSubstitutedFirstHalf(encoded));  
                Assert.AreEqual((uint) test.sumOfSkills, serial.getSumOfSkills(encoded));  
                Assert.AreEqual((uint) test.generation, serial.getGeneration(encoded));  
                Assert.AreEqual((bool) test.outOfGameFirstHalf, serial.getOutOfGameFirstHalf(encoded));  
                Assert.AreEqual((bool) test.yellowCardFirstHalf, serial.getYellowCardFirstHalf(encoded));  
            }
        }  

        [TestMethod]
        public void DecodeTactics() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingTacticsTestData.json");
            foreach(dynamic test in tests) {
                BigInteger encoded;
                bool succeeded = BigInteger.TryParse((string) test.encodedTactics, out encoded);
                Assert.AreEqual(true, succeeded);  
                Assert.AreEqual((uint) test.tacticsId, serial.getTacticsId(encoded));  
                for (int i = 0; i < 10; i++) { Assert.AreEqual((bool) test.extraAttack[i], serial.getExtraAttack(encoded, i)); }
                for (int i = 0; i < 3; i++) { Assert.AreEqual((uint) test.substitution[i], serial.getSubstitution(encoded, i)); }
                for (int i = 0; i < 3; i++) { Assert.AreEqual((uint) test.subsRound[i], serial.getSubsRound(encoded, i)); }
                for (int i = 0; i < 14; i++) { Assert.AreEqual((uint) test.linedUp[i], serial.getLinedUp(encoded, i)); }
            }
        }  

        [TestMethod]
        public void EncodeTactics() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingTacticsTestData.json");
            foreach(dynamic test in tests) {
                BigInteger encodedExpected;
                bool succeeded = BigInteger.TryParse((string) test.encodedTactics, out encodedExpected);
                Assert.AreEqual(true, succeeded); 
                uint[] substitution = tu.DynamicToUintArray(test.substitution);
                uint[] subsRound = tu.DynamicToUintArray(test.subsRound);
                uint[] linedUp = tu.DynamicToUintArray(test.linedUp);
                bool[] extraAttack = tu.DynamicToBoolArray(test.extraAttack);
                (BigInteger encoded, string err) = serial.encodeTactics(
                    substitution,
                    subsRound,
                    linedUp,
                    extraAttack,
                    (uint) test.tacticsId
                );
                Assert.AreEqual(err, "");  
                Assert.AreEqual(encoded, encodedExpected);  
            }
        }  

        [TestMethod]
        public void DecodeMatchLog() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingMatchLogTestData.json");
            foreach(dynamic test in tests) {
                BigInteger encoded;
                bool succeeded = BigInteger.TryParse((string) test.encodedLog, out encoded);
                Assert.AreEqual(true, succeeded);  
                Assert.AreEqual((bool) test.isHomeStadium, serial.getIsHomeStadium(encoded));  
                Assert.AreEqual((uint) test.winner, serial.getWinner(encoded));  
                Assert.AreEqual((uint) test.teamSumSkills, serial.getTeamSumSkills(encoded));  
                Assert.AreEqual((uint) test.trainingPoints, serial.getTrainingPoints(encoded));  
                Assert.AreEqual((uint) test.nGoals, serial.getNGoals(encoded));  
                Assert.AreEqual((uint) test.changesAtHalftime, serial.getChangesAtHalfTime(encoded));  
                int MAX_GOALS = 12;
                int N_HALFS = 2;
                for (int i = 0; i < MAX_GOALS; i++) { Assert.AreEqual((uint) test.assister[i], serial.getAssister(encoded, i)); }
                for (int i = 0; i < MAX_GOALS; i++) { Assert.AreEqual((uint) test.shooter[i], serial.getShooter(encoded, i)); }
                for (int i = 0; i < MAX_GOALS; i++) { Assert.AreEqual((uint) test.forwardPos[i], serial.getForwardPos(encoded, i)); }
                for (int i = 0; i < 7; i++) { Assert.AreEqual((bool) test.penalty[i], serial.getPenalty(encoded, i)); }
                for (int i = 0; i < 3; i++) { Assert.AreEqual((uint) test.halfTimeSubs[i], serial.getHalfTimeSubs(encoded, i)); }
                for (int h = 0; h < N_HALFS; h++) { Assert.AreEqual((uint) test.nGKAndDefs[h], serial.getNGKAndDefs(encoded, h == 1)); }
                for (int h = 0; h < N_HALFS; h++) { Assert.AreEqual((uint) test.nTot[h], serial.getNTot(encoded, h == 1)); }
                for (int h = 0; h < N_HALFS; h++) { Assert.AreEqual((uint) test.outOfGamePlayer[h], serial.getOutOfGamePlayer(encoded, h == 1)); }
                for (int h = 0; h < N_HALFS; h++) { Assert.AreEqual((uint) test.outOfGameType[h], serial.getOutOfGameType(encoded, h == 1)); }
                for (int h = 0; h < N_HALFS; h++) { Assert.AreEqual((uint) test.outOfGameRound[h], serial.getOutOfGameRound(encoded, h == 1)); }
                int counter = 0;
                for (int half = 0; half < N_HALFS; half++) {
                    for (int posInHalf = 0; posInHalf < 2; posInHalf++) {
                        Assert.AreEqual((uint) test.yellowCard[counter], serial.getYellowCard(encoded, posInHalf, half == 1));
                        counter++;
                    }
                }
                counter = 0;
                for (int half = 0; half < N_HALFS; half++) {
                    for (int posInHalf = 0; posInHalf < 3; posInHalf++) {
                        Assert.AreEqual((uint) test.inGameSubsHappened[counter], serial.getInGameSubsHappened(encoded, posInHalf, half == 1));
                        counter++;
                    }
                }

            }
        } 

        [TestMethod]
        public void hashes() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            Assert.AreEqual((ulong) 4623503348185510199, serial.int_hash("hola")); 
            BigInteger big = new BigInteger(534298574);
            Assert.AreEqual((ulong) 4623503348185510199, serial.int_hash(big.ToString())); 
        }   
    }
}
