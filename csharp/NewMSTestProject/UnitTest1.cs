using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Numerics;
using System;

namespace NewMSTestProject
{
    [TestClass]
    public class UnitTest1
    {
        [TestMethod]
        public void Test_AddMethod() {  
            Serialization serial = new Serialization();  
            double res = serial.Add(10, 20);  
            Assert.AreEqual(res, 30);  
        }  
        [TestMethod]
        public void Test_AddBNMethod() {  
            Serialization serial = new Serialization();  
            BigInteger res = serial.AddBN(new BigInteger(10), new BigInteger(20));  
            Assert.AreEqual(res, new BigInteger(30));  
        }  

        [TestMethod]
        public void getCurrentShirtNum() {  
            Serialization serial = new Serialization();
            uint shirt = 13;
            BigInteger state = new BigInteger(shirt * Math.Pow(2,43));
            uint res = serial.getCurrentShirtNum(state);  
            Assert.AreEqual(res, (uint) shirt);  
        }  

        [TestMethod]
        public void Test_encodingSkills() {  
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
        public void Test_encodingTactics() {  
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
        public void Test_encodingMatchLog() {  
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
                // for (int i = 0; i < 10; i++) { Assert.AreEqual((bool) test.extraAttack[i], serial.getExtraAttack(encoded, i)); }
                // for (int i = 0; i < 3; i++) { Assert.AreEqual((uint) test.substitution[i], serial.getSubstitution(encoded, i)); }
                // for (int i = 0; i < 3; i++) { Assert.AreEqual((uint) test.subsRound[i], serial.getSubsRound(encoded, i)); }
                // for (int i = 0; i < 14; i++) { Assert.AreEqual((uint) test.linedUp[i], serial.getLinedUp(encoded, i)); }
            }
        }  
        
    }
}
