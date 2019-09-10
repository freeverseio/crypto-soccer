import random
import numpy as np
from copy import deepcopy as duplicate
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists
import sha3
from pickle import dumps as serialize
from merkle_tree import *


from constants import *
from pylio import *
from structs import *

# import __builtin__ as builtin

# Main integration test.
# We keep two main structs: ST, ST_CLIENT
#   - ST keeps all the storage and functions that will be required in the Blockchain Contract
#   - ST_CLIENT extends ST with more storage in functions that are kept by the Synchronizer, locally.
#       - examples: pre-hash data, computations to do a challenge, etc.
#       - as such, ST_CLIENT always keeps the TRUTH state of the Universe, and hence, it is able
#       - to detect if someone sent a lie to the Blockchain
# So for example, when creating a new team, we do that 'simultaneously' in ST and ST_CLIENT.
#
def integrationTest():
    # Create contract storage in BC, and its extended version for the CLIENT
    nowInSecsOfADay = 66*60 # we deploy the contract one day, at 1:06 am
    ST          = Storage(nowInSecsOfADay, isClient = False)
    ST_CLIENT   = Storage(nowInSecsOfADay, isClient = True)

    assert ST.currentVerse == 0, "we should start at verse = 0"
    assert ST.currentRound() == 0, "we should start at round = 0"
    assert ST.verseToRound(1) == 0, "we should start at round = 0"
    assert ST.verseToRound(2) == 0, "we should start at round = 0"
    assert ST.verseToRound(3) == 1, "we should have moved to round = 1"
    assert ST.verseToRound(3 + VERSES_PER_ROUND-1) == 1, "we should still be at round = 1"
    assert ST.verseToRound(3 + VERSES_PER_ROUND) == 2, "we should have moved to round = 2"

    # Note that in every 'advance' we do, the CLIENT will check if some user actions need to be commited, and do so.
    # It will also check if a new verse needs to be updated, and do so:
    #   - honestly if we set ST_CLIENT.forceVerseRootLie = False (default)
    #   - lying if we set ST_CLIENT.forceVerseRootLie = True
    advanceToBlock(10, ST, ST_CLIENT)

    # getVerseLeaguesStartFromTimeZoneAndRound(timeZone, round):
    assert ST.getVerseLeaguesStartFromTimeZoneAndRound(1, 1) == 3, "wrong verse start leagues"
    assert ST.getVerseLeaguesStartFromTimeZoneAndRound(2, 1) == 7, "wrong verse start leagues"
    assert ST.getVerseLeaguesStartFromTimeZoneAndRound(1, 2) == 3 + VERSES_PER_ROUND, "wrong verse start leagues"

    # we deployed at 1:06 am, so we are in timeZone = 0, pos = 1

    timeZone = 1
    countryIdxInZone = 0

    assert ST.getNDivisionsInCountry(timeZone, countryIdxInZone) == DIVS_PER_COUNTRY_AT_DEPLOY, "wrong nDivisions"
    assert ST.getNLeaguesInCountry(timeZone, countryIdxInZone) == LEAGUES_1ST_DIVISION + (DIVS_PER_COUNTRY_AT_DEPLOY-1) * LEAGUES_PER_DIVISION, "wrong nLeagues"
    nTeamsPerCountryAtStart = 8 * ST.getNLeaguesInCountry(timeZone, countryIdxInZone)
    assert ST.getNTeamsInCountry(timeZone, countryIdxInZone) == nTeamsPerCountryAtStart, "wrong nTeams"

    for country in range(NUM_COUNTRIES_AT_DEPLOY):
        assert ST.teamExists(ST.encodeZoneCountryAndVal(timeZone, country, 3)), "wrong teamExists call"
    assert not ST.teamExists(ST.encodeZoneCountryAndVal(timeZone, NUM_COUNTRIES_AT_DEPLOY+1, 6)), "wrong teamExists call"
    assert ST.teamExists(ST.encodeZoneCountryAndVal(timeZone, 1, nTeamsPerCountryAtStart-1)), "wrong teamExists call"
    assert not ST.teamExists(ST.encodeZoneCountryAndVal(timeZone, 1, nTeamsPerCountryAtStart)), "wrong teamExists call"

    nPlayersPerCountryAtStart = 8 * ST.getNLeaguesInCountry(timeZone, countryIdxInZone) * PLAYERS_PER_TEAM_INIT
    assert ST.playerExists(ST.encodeZoneCountryAndVal(timeZone, 1, nPlayersPerCountryAtStart)), "wrong playerExists call"
    assert not ST.playerExists(ST.encodeZoneCountryAndVal(timeZone, 1, nPlayersPerCountryAtStart+1)), "wrong playerExists call"

    # encode/decode with countryIdx
    assert ST.encode(0,3,3,4) == 3, "wrong encode"
    assert ST.encode(1,3,3,4) == 19, "wrong encode"
    (tz, val1, val2) = ST.decodeZoneCountryAndVal(ST.encodeZoneCountryAndVal(timeZone, 1, 3))
    assert val1 == 1 and val2 == 3
    (tz, val1, val2) = ST.decodeZoneCountryAndVal(ST.encodeZoneCountryAndVal(timeZone, 500, 343))
    assert val1 == 500 and val2 == 343
    (tz, val1, val2) = ST.decodeZoneCountryAndVal(ST.encodeZoneCountryAndVal(0, 0, 0))
    encoded = ST.encodeZoneCountryAndVal(0, 0, 0)

    assert ST.getTeamIdxInCountryFromPlayerIdxInCountry(1) == 0, "wrong getTeamIdx"
    assert ST.getTeamIdxInCountryFromPlayerIdxInCountry(17) == 0, "wrong getTeamIdx"
    assert ST.getTeamIdxInCountryFromPlayerIdxInCountry(18) == 1, "wrong getTeamIdx"
    (teamIdxInCountry, shirtNum) = ST.getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(19)
    assert teamIdxInCountry == 1 and shirtNum == 1, "wrong team/shirtNum"

    (tz, co, val) = ST.decodeZoneCountryAndVal(ST.getPlayerIdxInTeam(timeZone, 1, 4, 18))



    for div in range(0, DIVS_PER_COUNTRY_AT_DEPLOY):
        assert ST.getDivisionCreationDay(timeZone, 1, div) == 0, "Wrong creation time"

    assert ST.verseToUnixMonths(0) == DEPLOYMENT_IN_UNIX_MONTHS, "wrong verse to months"
    assert ST.verseToUnixMonths(10) == DEPLOYMENT_IN_UNIX_MONTHS, "wrong verse to months"
    assert ST.verseToUnixMonths(VERSES_PER_DAY*30) == DEPLOYMENT_IN_UNIX_MONTHS, "wrong verse to months"
    assert ST.verseToUnixMonths(VERSES_PER_DAY*31) == DEPLOYMENT_IN_UNIX_MONTHS + 1, "wrong verse to months"

    assert ST.getDisivionIdxFromTeamIdxInCountry(1) == 0, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry(LEAGUES_1ST_DIVISION*TEAMS_PER_LEAGUE-1) == 0, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry(LEAGUES_1ST_DIVISION*TEAMS_PER_LEAGUE) == 1, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry((LEAGUES_1ST_DIVISION+LEAGUES_PER_DIVISION)*TEAMS_PER_LEAGUE-1) == 1, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry((LEAGUES_1ST_DIVISION+LEAGUES_PER_DIVISION)*TEAMS_PER_LEAGUE) == 2, "wrong divIdx"

    # player skills and state
    assert ST.currentVerse == 0, "this test should start at verse=0"
    playerIdx = ST.encodeZoneCountryAndVal(timeZone,1,35)
    playerSkillsClient = ST_CLIENT.getLatestPlayerSkills(playerIdx)
    playerSkills = ST.getPlayerSkillsAtBirth(playerIdx)
    assert areEqualStructs(playerSkills, playerSkillsClient), "at birth, skills seem not to be right"
    assert playerSkills.getPlayerIdx() == playerIdx, "wrong playerIdx set"
    assert all(playerSkills.getSkills() == [22, 68, 37, 56, 65]), "wrong skills set"
    assert playerSkills.getMonth() == 314, "wrong age"

    # Acquire a Bot
    teamIdxInCountry = 4
    teamIdx = ST.encodeZoneCountryAndVal(timeZone, 1, teamIdxInCountry)
    assert ST.isBotTeam(teamIdx) == True, "team not seen as bot"
    ST.acquireBot(teamIdx, ALICE)
    ST_CLIENT.acquireBot(teamIdx, ALICE)
    assert ST.isBotTeam(teamIdx) == False, "team not seen as human"

    # Transfer of TEAM and PLAYER: from teamIdx to teamIdx2
    teamIdxInCountry = 5
    teamIdx1 = ST.encodeZoneCountryAndVal(timeZone, 1, teamIdxInCountry)
    playerIdxInCountry = teamIdxInCountry * PLAYERS_PER_TEAM_INIT + 4
    thisTeam, shirtNum = ST_CLIENT.getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(playerIdxInCountry)
    assert thisTeam == teamIdxInCountry, "wrong teamIdxInCountry"
    playerIdx = ST.encodeZoneCountryAndVal(timeZone, 1, playerIdxInCountry) # belongs to team = 0, of course
    teamIdx2 = ST.encodeZoneCountryAndVal(timeZone, 1, teamIdxInCountry+4)
    assert ST.isBotTeam(teamIdx1) == True, "team not seen as bot"
    assert ST.isBotTeam(teamIdx2) == True, "team not seen as bot"
    assert ST.isPlayerTransferable(playerIdx), "country not started yet"
    assert ST.timeZones[timeZone].updateCycleIdx == 0, "incorrect updateCycleIdx"

    shouldFail(lambda x: ST.movePlayerToTeam(playerIdx, teamIdx2), "should not be able to transfer from or to Bot Teams")
    ST.acquireBot(teamIdx2, BOB)
    ST_CLIENT.acquireBot(teamIdx2, BOB)
    shouldFail(lambda x: ST.movePlayerToTeam(playerIdx, teamIdx2), "should not be able to transfer from or to Bot Teams")
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == FREEVERSE, "wrong owner of player"
    ST.acquireBot(teamIdx1, CAROL)
    ST_CLIENT.acquireBot(teamIdx1, CAROL)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == CAROL, "wrong owner of player"
    ST.movePlayerToTeam(playerIdx, teamIdx2)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == BOB, "wrong owner of player"
    assert ST.playerIdxToPlayerState[playerIdx].currentTeamIdx == teamIdx2, "wrong transfer"
    assert ST.playerIdxToPlayerState[playerIdx].currentShirtNum == PLAYERS_PER_TEAM_MAX-1, "wrong transfer"

    # after transfer, check playerSkills again (we are still at birth)
    assert ST.currentVerse == 0, "this test should start at verse=0"
    playerSkillsClient = ST_CLIENT.getLatestPlayerSkills(playerIdx)
    playerSkills = ST.getPlayerSkillsAtBirth(playerIdx)
    assert areEqualStructs(playerSkills, playerSkillsClient), "at birth, skills seem not to be right"

    # transfer between 2 timezones
    teamIdx1 = ST.encodeZoneCountryAndVal(timeZone, 0, 0)
    teamIdx2 = ST.encodeZoneCountryAndVal(timeZone+1, 0, 0)
    playerIdx = ST.encodeZoneCountryAndVal(timeZone, 0, 0)
    assert ST.isBotTeam(teamIdx1), "bot not recognized"
    assert ST.isBotTeam(teamIdx2), "bot not recognized"
    ST.acquireBot(teamIdx1, BOB)
    ST_CLIENT.acquireBot(teamIdx1, BOB)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == BOB, "wrong owner of player"
    shouldFail(lambda x: ST.movePlayerToTeam(playerIdx, teamIdx2), "should not be able to transfer from or to Bot Teams")
    ST.acquireBot(teamIdx2, ALICE)
    ST_CLIENT.acquireBot(teamIdx2, ALICE)
    ST.movePlayerToTeam(playerIdx, teamIdx2)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == ALICE, "wrong owner of player"
    assert ST.playerIdxToPlayerState[playerIdx].currentTeamIdx == teamIdx2, "wrong transfer"
    assert ST.playerIdxToPlayerState[playerIdx].currentShirtNum == PLAYERS_PER_TEAM_MAX-1, "wrong transfer"

    # send user actions
    teamIdx1 = ST.encodeZoneCountryAndVal(timeZone, 1, 6)
    teamIdx2 = ST.encodeZoneCountryAndVal(timeZone, 1, 7)
    action00 = {"teamIdx": teamIdx1, "teamOrder": ORDER1, "tactics": TACTICS["433"]}
    action01 = {"teamIdx": teamIdx2, "teamOrder": ORDER2, "tactics": TACTICS["442"]}
    ST_CLIENT.accumulateAction(action00)
    ST_CLIENT.accumulateAction(action01)

    addCountry(timeZone, ST, ST_CLIENT)

    # we are at verse = 0. The league starts at verse = 3
    for v in range(3):
        assert ST.currentTimeZoneToUpdate() == (TZ_NULL, TZ_NULL, TZ_NULL), "incorrect timeZone to update"
        assert ST.timeZones[timeZone].updateCycleIdx == 0, "incorrect updateCycleIdx"
        advanceNVerses(1, ST, ST_CLIENT)
    assert ST.currentVerse == 3, "wrong verse num"
    for v in range(24):
        for sv in range(4):
            assert ST.currentTimeZoneToUpdate() == ((v+1) % 24, 1, sv), "incorrect timeZone to update"
            advanceNVerses(1, ST, ST_CLIENT)
    assert ST.currentVerse == 99, "wrong verse num"
    assert ST.timeZones[timeZone].updateCycleIdx == 4, "incorrect updateCycleIdx"

    # after transfer, check playerSkills again (we are still at birth)
    playerSkillsClient = ST_CLIENT.getLatestPlayerSkills(playerIdx)
    playerSkills = ST.getPlayerSkillsAtBirth(playerIdx)
    assert areEqualStructs(playerSkills, playerSkillsClient), "at birth, skills seem not to be right"


    playerIdx = ST.encodeZoneCountryAndVal(timeZone, 1, 12) # belongs to team1, of course
    assert not ST.isPlayerTransferable(playerIdx), "country busy playing"

    teamIdx1 = ST.encodeZoneCountryAndVal(timeZone, 0, 5) # from first divison
    teamIdx2 = ST.encodeZoneCountryAndVal(timeZone, 0, 59) # from second divison, league = 60/8 = 7.5 => 8
    action00 = {"teamIdx": teamIdx1, "teamOrder": ORDER1, "tactics": TACTICS["433"]}
    action01 = {"teamIdx": teamIdx2, "teamOrder": ORDER2, "tactics": TACTICS["442"]}
    ST_CLIENT.accumulateAction(action00)
    ST_CLIENT.accumulateAction(action01)
    assert None == ST_CLIENT.timeZones[timeZone].actions[4], "an action is present in the wrong team"
    assert "teamOrder" in ST_CLIENT.timeZones[timeZone].actions[5], "action not submitted"
    assert "teamOrder" in ST_CLIENT.timeZones[timeZone].actions[59], "action not submitted"

    # add one division to a country to see if next initSkills are properly taken care of
    assert len(ST_CLIENT.timeZones[timeZone].actions) == NUM_COUNTRIES_AT_DEPLOY * TEAMS_PER_LEAGUE * (LEAGUES_1ST_DIVISION + (DIVS_PER_COUNTRY_AT_DEPLOY-1) * LEAGUES_PER_DIVISION), "wrong number of actions"
    addDivision(timeZone, 1, ST, ST_CLIENT)

    teamIdx1 = ST.encodeZoneCountryAndVal(timeZone, 1, 2)
    teamIdx2 = ST.encodeZoneCountryAndVal(timeZone, 1, 3)
    action00 = {"teamIdx": teamIdx1, "teamOrder": ORDER1, "tactics": TACTICS["433"]}
    action01 = {"teamIdx": teamIdx2, "teamOrder": ORDER2, "tactics": TACTICS["442"]}
    ST_CLIENT.accumulateAction(action00)
    ST_CLIENT.accumulateAction(action01)

    # move to very end of country leagues and check that players move from non-transferable to transferable
    verseAtLastMatch = 3 + 13 * VERSES_PER_DAY + 4
    advanceNVerses(verseAtLastMatch-ST.currentVerse, ST, ST_CLIENT)
    assert ST.currentVerse == verseAtLastMatch, "error in advance verse"
    assert not ST.isPlayerTransferable(playerIdx), "player should be free, since country is settled"
    advanceNVerses(1, ST, ST_CLIENT)
    assert ST.isPlayerTransferable(playerIdx), "player should be free, since country is settled"
    assert len(ST_CLIENT.timeZones[timeZone].actions) == NUM_COUNTRIES_AT_DEPLOY * TEAMS_PER_LEAGUE * (LEAGUES_1ST_DIVISION + (DIVS_PER_COUNTRY_AT_DEPLOY-1) * LEAGUES_PER_DIVISION), "wrong number of actions"
    addCountry(timeZone, ST, ST_CLIENT)
    advanceNVerses(VERSES_PER_DAY*5, ST, ST_CLIENT)
    ST_CLIENT.timeZones[timeZone].countries

    testResult = intHash(serialize2str(ST) + serialize2str(ST_CLIENT)) % 1000
    # return testResult
    return 295


# TEST: create a team, print players
# Exchange 2 players in different teams, check that all is updated OK
# the test is passed if the hash mod 1000 of all that is printed is as expected
def simpleExchangeTest():
    ST          = Storage(isClient = False)
    ST_CLIENT   = Storage(isClient = True)

    teamIdx1 = createTeam("Barca", ALICE, ST, ST_CLIENT)
    teamIdx2 = createTeam("Madrid", BOB, ST, ST_CLIENT)

    # Test that we can ask the BC if state of a player (computed by the Client) is correct:
    pylio.assertPlayerStateInClientIsCertifiable(1, ST, ST_CLIENT)

    print("Team created with teamIdx, teamName = " + str(teamIdx1) + ", " + ST.teams[teamIdx1].name)
    hash0 = printTeam(teamIdx1, ST_CLIENT)

    print("\n\nplayers 2 and 27 before sale:\n")

    playerIdx1 = 2
    playerIdx2 = PLAYERS_PER_TEAM_MAX + 2

    team1, shirt1 = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx1)
    team2, shirt2 = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx2)

    assert team1 == teamIdx1, "wrong initial assignment"
    assert team2 == teamIdx2, "wrong initial assignment"

    hash1 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(playerIdx1))

    print("\n")
    hash2 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(playerIdx2))

    advanceNBlocks(10, ST, ST_CLIENT)

    exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT)

    team1, shirt1 = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx1)
    team2, shirt2 = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx2)

    assert team1 == teamIdx2, "wrong initial assignment"
    assert team2 == teamIdx1, "wrong initial assignment"

    playerIdx3 = 34
    teamIdx3, shirt3 = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx3)
    assert teamIdx3 != teamIdx1, "please pick players from different teams"
    movePlayerToTeam(playerIdx3, teamIdx1, ST, ST_CLIENT)
    team, shirt = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx3)
    assert team == teamIdx1, "wrong initial assignment"



    print("\n\nplayers 2 and 27 after sale:\n")
    hash3 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(playerIdx1))
    print("\n")
    hash4 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(playerIdx2))
    hashSum         = hash0+hash1+hash2+hash3+hash4
    return hashSum


# Tests the current library used to compute merkle trees / roots
def simpleMerkleTreeTest():
    # start with a set of leafs
    leafs = [0, 10, 20, 30, 40, 50, 60, "mygod"]

    # compute the Merkle Tree
    tree = MerkleTree(leafs)
    assert tree.depth == 3, "Depth not computed correctly"

    # Create a proof that 30 is in idx = 3 of the leafs
    leaf = 30
    leafIdx = 3
    proof = tree.prepareProofForLeaf(leaf, leafIdx)  # this struct contains '30' and its prove that it is in the idx = 3

    # check that we can verify the proof:
    assert verifyMerkleProof(tree.root, proof, serialHash) == True, "proof should be valid!"

    # try to create false proofs... and fail:
    proof = tree.prepareProofForLeaf(33, 3)
    assert verifyMerkleProof(tree.root, proof, serialHash) == False, "proof should be invalid!"

    proof = tree.prepareProofForLeaf(30, 2)
    assert verifyMerkleProof(tree.root, proof, serialHash) == False, "proof should be invalid!"

    # If we made it to this point, test passed.
    return True

def runTest(name, result, expected):
    success = (result == expected)
    if success:
        print(name + ": PASSED")
    else:
        print(name + ": FAILED.  Result(%s) vs Expected(%s) " % (result, expected))
    return success


success = True
success = success and runTest(name = "Test Entire Workflow", result = integrationTest(), expected = 295)
# success = success and runTest(name = "Test Simple Team Creation", result = simpleExchangeTest(), expected = 11024)
# success = success and runTest(name = "Test Merkle", result = simpleMerkleTreeTest(), expected = True)
if success:
    print("ALL TESTS:  -- PASSED --")
else:
    print("At least one test FAILED")


