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
        public void readJson() {  
            Serialization serial = new Serialization();
            TestUtils tu = new TestUtils();
            dynamic tests = tu.LoadJson("encodingSkillsTestData.json");
            foreach(dynamic test in tests) {
                BigInteger encodedSkills;
                bool succeeded = BigInteger.TryParse((string) test.encodedSkills, out encodedSkills);
                Assert.AreEqual(true, succeeded);  
                for (int sk = 0; sk < 5; sk++) { Assert.AreEqual((uint) test.skills[sk], serial.getSkill(encodedSkills, sk)); }
                Assert.AreEqual((uint) test.birthday, serial.getBirthDay(encodedSkills));  
                Assert.AreEqual((bool) test.isSpecial, serial.getIsSpecial(encodedSkills));  
                Assert.AreEqual((uint) test.potential, serial.getPotential(encodedSkills));  
                Assert.AreEqual((uint) test.forwardness, serial.getForwardness(encodedSkills));  
                Assert.AreEqual((uint) test.leftishness, serial.getLeftishness(encodedSkills));  
                Assert.AreEqual((uint) test.aggressiveness, serial.getAggressiveness(encodedSkills));  
                Assert.AreEqual((bool) test.alignedEndOfFirstHalf, serial.getAlignedEndOfFirstHalf(encodedSkills));  
                Assert.AreEqual((bool) test.redCardLastGame, serial.getRedCardLastGame(encodedSkills));  
                Assert.AreEqual((uint) test.gamesNonStopping, serial.getGamesNonStopping(encodedSkills));  
                Assert.AreEqual((uint) test.injuryWeeksLeft, serial.getInjuryWeeksLeft(encodedSkills));  
                Assert.AreEqual((bool) test.substitutedFirstHalf, serial.getSubstitutedFirstHalf(encodedSkills));  
                Assert.AreEqual((uint) test.sumOfSkills, serial.getSumOfSkills(encodedSkills));  
                Assert.AreEqual((uint) test.generation, serial.getGeneration(encodedSkills));  
                Assert.AreEqual((bool) test.outOfGameFirstHalf, serial.getOutOfGameFirstHalf(encodedSkills));  
                Assert.AreEqual((bool) test.yellowCardFirstHalf, serial.getYellowCardFirstHalf(encodedSkills));  
                

            }
        }  
    }
}
