using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Numerics;
using System;
using static Serialization;

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
            // This test vector comes from Golang tests.
            BigInteger big = new BigInteger(123456789);
            uint unsig = 123456789;
            string[] inputs = new string[3]{"hola", big.ToString(), unsig.ToString()};
            ulong[] expectedOutputs = new ulong[3]{4623503348185510199, 492395637191921148, 492395637191921148};
            for (uint i = 0; i < inputs.Length; i++) {
                Assert.AreEqual(expectedOutputs[i], serial.int_hash(inputs[i])); 
            }
        }   

        [TestMethod]
        public void matchEvents() {  
            // From TestMatchEvents2ndHalfHardcoded
            Serialization serial = new Serialization();
            BigInteger verseSeed = new BigInteger(new byte[32]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2});
            Console.WriteLine(verseSeed);
            BigInteger[] teamIds = new BigInteger[2]{1,2};
            BigInteger log = BigInteger.Parse("452312848584470512245079946786433186608365459112320500501947696564481818624");
            BigInteger[] matchLogsAndEvents = new BigInteger[2 + 5 * 12] { 
                log, log,
                1, 0, 0, 0, 0,
                1, 0, 0, 0, 0,
                0, 1, 7, 1, 7,
                1, 0, 0, 0, 0,
                0, 0, 0, 0, 0,
                1, 0, 0, 0, 0,
                1, 0, 0, 0, 0,
                0, 1, 10, 1, 10,
                0, 1, 7, 1, 7,
                0, 0, 0, 0, 0,
                1, 0, 0, 0, 0,
                0, 0, 0, 0, 0,
            };
            uint NO_SUBS = 11;
            uint NO_LINEUP = 25;
            BigInteger tact1, tact2;
            string err;
            (tact1, err) = serial.encodeTactics(
                new uint[3]{5, 1, NO_SUBS},
                new uint[3]{4, 6, 7},
                new uint[14]{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, NO_LINEUP},
                new bool[10]{false, false, false, false, false, false, false, false, false, false},
                0
            );         
            Assert.AreEqual("", err);
            (tact2, err) = serial.encodeTactics(
                new uint[3]{5, 1, NO_SUBS},
                new uint[3]{4, 6, 7},
                new uint[14]{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, NO_LINEUP},
                new bool[10]{false, false, false, false, false, false, false, false, false, false},
                0
            );
            Assert.AreEqual("", err);
            BigInteger[] tactics = new BigInteger[2]{tact1, tact2};
            BigInteger[] ids = new BigInteger[25];
            for (uint p = 0; p < ids.Length; p++) { ids[p] = 20000+p; } // any playerId works as long as it is > 2.
            BigInteger[][] playerIds = new BigInteger[2][]{ids, ids};
            bool is2ndHalf = true;
            MatchEvent[] events;
            (events, err) = serial.processMatchEvents(is2ndHalf, verseSeed, teamIds, tactics, playerIds, matchLogsAndEvents);
            string concat = "";
            uint[] nGoals = new uint[2]{0, 0};
            for (uint i = 0; i < events.Length; i++) {
                concat += "[";
                concat += events[i].Minute.ToString();
                concat += ", ";
                concat += events[i].Type.ToString();
                concat += ", ";
                concat += events[i].Team.ToString();
                concat += ", ";
                concat += (events[i].ManagesToShoot ? "true" : "false");
                concat += ", ";
                concat +=(events[i].IsGoal ? "true" : "false");
                concat += ", ";
                if (events[i].IsGoal) {
                    nGoals[events[i].Team]++;
                }
                concat += events[i].PrimaryPlayer.ToString();
                concat += ", ";
                concat += events[i].SecondaryPlayer.ToString();
                concat += "]";
            }
            Console.WriteLine(concat);
            string expectedConcat = "[46, 0, 1, false, false, 16, -1][52, 0, 1, false, false, 0, -1][55, 0, 0, true, true, 8, 8][58, 0, 1, false, false, 15, -1][61, 0, 0, false, false, 1, -1][68, 0, 1, false, false, 13, -1][71, 0, 1, false, false, 8, -1][74, 0, 0, true, true, 10, 10][77, 0, 0, true, true, 8, 8][84, 0, 0, false, false, 4, -1][86, 0, 1, false, false, 11, -1][91, 0, 0, false, false, 2, -1][68, 2, 0, false, false, 12, -1][47, 1, 0, false, false, 13, -1][68, 2, 1, false, false, 18, -1][78, 1, 1, false, false, 0, -1][61, 5, 0, false, false, 11, 19][67, 5, 0, false, false, 16, 12][61, 5, 1, false, false, 1, 17][67, 5, 1, false, false, 4, 18]";
            Assert.AreEqual(concat, expectedConcat);
        }   
    }
}
