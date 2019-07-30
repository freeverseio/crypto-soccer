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

# Function used in the integration test to update a set of leagues always telling the truth
def updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, doExchanges):
    for extraVerse in range(2000):
        if doExchanges and extraVerse % 10:
            for p in range(2):
                playerIdx1 = 1 + intHash(str(p+extraVerse)) % 100 * PLAYERS_PER_TEAM_MAX
                playerIdx2 = 1 + intHash(str(p+extraVerse) + "salt") % 100 * PLAYERS_PER_TEAM_MAX
                try:
                    exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT)
                except:
                    pass
                pylio.assertPlayerStateInClientIsCertifiable(playerIdx1, ST, ST_CLIENT)

        for leagueIdx in leaguesTested:
            verse = ST.leagues[leagueIdx].verseFinal()
            verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
            assert not ( ST.isLeagueSettled(leagueIdx) and verseStatus != UPDT_LEVEL1), "Someone hacked the game"
            if ST.hasLeagueFinished(leagueIdx) and (not ST.isLeagueSettled(leagueIdx)):
                print("challenging league...", leagueIdx)
                if verseStatus == UPDT_LEVEL2:
                    print("challenging league... superRoot", leagueIdx)
                    subVerse = 0
                    challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie=0)
                    ST.assertCanChallengeStatus(verse, UPDT_LEVEL3)
                if verseStatus == UPDT_LEVEL3:
                    print("challenging league... leagueRoot", leagueIdx)
                    challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie=0)
                elif verseStatus == UPDT_LEVEL4:
                    thisLeagueIdx = ST.getLeagueIdxFromPosInSubverse(verse, ST.verseToLeagueCommits[verse].posInSubVerse)
                    print("challenging league... initSkills", thisLeagueIdx)
                    challengeLevel4(LEAGUE_INIT_SKILLS_ID, verse, ST, ST_CLIENT)
        return ST, ST_CLIENT

# Function to create a set of updates/challenges to many leagues
# We start with all such leagues already updated with the "VerseRoot", which was always TRUE/HONEST => Level 1 is True
# Then:
#   - if we encounter the league in Level 1 / VerseRoot  => challenge with a Lie (moving to Level 2 / SuperRoots)
#   - if we encounter the league in Level 2 / SuperRoots => challenge with a Lie (moving to Level 3 / LeagueRoots)
#   - if we encounter the league in Level 3 / LeagueRoots:
#       - if first time for this league  => challenge with a Lie (moving to Level 4 / OneLeague)
#       - if second time for this league => challenge with a Truth (moving to Level 4 / OneLeague)
#   - if we encounter the league in Level 4 / OneLeague:
#         - challenge it, and check that it is successful only if it had truly lied.
def brutalBlock(ST, ST_CLIENT, leaguesTested):
    leaguesTestedAtLevel3 = []
    advanceNVerses(250, ST, ST_CLIENT)
    for extraVerse in range(45):
        advanceNVerses(1, ST, ST_CLIENT)
        for leagueIdx in leaguesTested:
            if ST.hasLeagueFinished(leagueIdx) and (not ST.isLeagueSettled(leagueIdx)):
                print("challenging league...", leagueIdx)
                verse = ST.leagues[leagueIdx].verseFinal()
                verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
                if verseStatus == UPDT_LEVEL1:
                    print("challenging league... verseRoot", leagueIdx)
                    challengeLevel1(verse, BOB, ST, ST_CLIENT, lie = 11)
                    ST.assertCanChallengeStatus(verse, UPDT_LEVEL2)
                elif verseStatus == UPDT_LEVEL2:
                    print("challenging league... superRoot", leagueIdx)
                    subVerse = 0
                    challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie=12)
                    ST.assertCanChallengeStatus(verse, UPDT_LEVEL3)
                elif verseStatus == UPDT_LEVEL3:
                    if leagueIdx in leaguesTestedAtLevel3:
                        print("challenging league... allLeagues with truth: ", leagueIdx)
                        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie=0)
                        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4)
                    else:
                        print("challenging league... allLeagues with lie: ", leagueIdx)
                        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie = 1)
                        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4)

                elif verseStatus == UPDT_LEVEL4:
                    thisLeagueIdx = ST.getLeagueIdxFromPosInSubverse(verse, ST.verseToLeagueCommits[verse].posInSubVerse)
                    print("challenging league... initSkills", thisLeagueIdx)
                    challengeLevel4(LEAGUE_INIT_SKILLS_ID, verse, ST, ST_CLIENT)
                    if thisLeagueIdx in leaguesTestedAtLevel3:
                        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4)
                    else:
                        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3)
                    leaguesTestedAtLevel3.append(thisLeagueIdx)
    return ST, ST_CLIENT




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

    # The accumulator is responsible for receiving user actions and committing them in the correct verse.
    # It only lives in the CLIENT.
    ST_CLIENT.addAccumulator()

    # Note that in every 'advance' we do, the CLIENT will check if some user actions need to be commited, and do so.
    # It will also check if a new verse needs to be updated, and do so:
    #   - honestly if we set ST_CLIENT.forceVerseRootLie = False (default)
    #   - lying if we set ST_CLIENT.forceVerseRootLie = True
    advanceToBlock(10, ST, ST_CLIENT)

    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(0)
    # assert timeZone == 0 and posInTimeZone == 1, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(1 + 0*4)
    # assert timeZone == 0 and posInTimeZone == 2, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(3 + 0*4)
    # assert timeZone == 1 and posInTimeZone == 0, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(3 + 1*4)
    # assert timeZone == 2 and posInTimeZone == 0, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(3 + 22*4)
    # assert timeZone == 23 and posInTimeZone == 0, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(23*4)
    # assert timeZone == 23 and posInTimeZone == 1, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(23*4 + 2)
    # assert timeZone == 23 and posInTimeZone == 3, "wrong timeZone"
    # timeZone, posInTimeZone = ST.verseToTimeZoneToUpdate(23*4 + 3)
    # assert timeZone == 0 and posInTimeZone == 0, "wrong timeZone"

    # getVerseLeaguesStartFromTimeZoneAndRound(timeZone, round):
    assert ST.getVerseLeaguesStartFromTimeZoneAndRound(1, 1) == 3, "wrong verse start leagues"
    assert ST.getVerseLeaguesStartFromTimeZoneAndRound(2, 1) == 7, "wrong verse start leagues"
    assert ST.getVerseLeaguesStartFromTimeZoneAndRound(1, 2) == 3 + VERSES_PER_ROUND, "wrong verse start leagues"

    # we deployed at 1:06 am, so we are in timeZone = 0, pos = 1
    # assert ST.currentTimeZoneToUpdate() == (0, 1), "wrong init timeZone"

    timeZone = 1
    countryIdx = 1

    assert ST.getNDivisionsInCountry(countryIdx) == 1, "wrong nDivisions"
    assert ST.getNLeaguesInCountry(countryIdx) == 1, "wrong nLeagues"
    assert ST.getNTeamsInCountry(countryIdx) == 8, "wrong nTeams"

    assert ST.teamExists(ST.encodeCountryAndVal(1, 3)), "wrong teamExists call"
    assert not ST.teamExists(ST.encodeCountryAndVal(2, 3)), "wrong teamExists call"
    assert ST.teamExists(ST.encodeCountryAndVal(1, 8)), "wrong teamExists call"
    assert not ST.teamExists(ST.encodeCountryAndVal(1, 9)), "wrong teamExists call"

    assert ST.playerExists(ST.encodeCountryAndVal(1, 3)), "wrong playerExists call"
    assert not ST.playerExists(ST.encodeCountryAndVal(2, 3)), "wrong playerExists call"
    assert ST.playerExists(ST.encodeCountryAndVal(1, 8*PLAYERS_PER_TEAM_INIT)), "wrong playerExists call"
    assert not ST.playerExists(ST.encodeCountryAndVal(1, 8*PLAYERS_PER_TEAM_INIT+1)), "wrong playerExists call"

    divisionIdx = addDivision(countryIdx, ST, ST_CLIENT)
    assert divisionIdx == 2, "wrong divisionIdx"
    assert ST.getNDivisionsInCountry(countryIdx) == 2, "wrong nDivisions"
    assert ST.getNLeaguesInCountry(countryIdx) == 17, "wrong nLeagues"
    assert ST.getNTeamsInCountry(countryIdx) == 17*8, "wrong nTeams"

    # getTeamIdxInCountryFromLeagueAndPos(divisionIdx, leaguePosInDiv, teamPosInLeague)
    shouldFail(lambda x: ST.getTeamIdxInCountryFromLeagueAndPos(0,1,1), "division 0 should not exist")
    shouldFail(lambda x: ST.getTeamIdxInCountryFromLeagueAndPos(1,1,0), "division 0 only has 1 league")
    assert ST.getTeamIdxInCountryFromLeagueAndPos(1, 0, 0) == 1, "wrong teamIdx"
    assert ST.getTeamIdxInCountryFromLeagueAndPos(1, 0, 1) == 2, "wrong teamIdx"
    assert ST.getTeamIdxInCountryFromLeagueAndPos(2, 0, 0) == 9, "wrong teamIdx"
    assert ST.getTeamIdxInCountryFromLeagueAndPos(2, 0, 1) == 10, "wrong teamIdx"
    assert ST.getTeamIdxInCountryFromLeagueAndPos(2, 1, 0) == 17, "wrong teamIdx"
    assert ST.getTeamIdxInCountryFromLeagueAndPos(2, 1, 1) == 18, "wrong teamIdx"

    # encode/decode with countryIdx
    assert ST.encode(0,3,3,4) == 3, "wrong encode"
    assert ST.encode(1,3,3,4) == 19, "wrong encode"
    (val1, val2) = ST.decodeCountryAndVal(ST.encodeCountryAndVal(1, 3))
    assert val1 == 1 and val2 == 3
    (val1, val2) = ST.decodeCountryAndVal(ST.encodeCountryAndVal(500, 343))
    assert val1 == 500 and val2 == 343

    assert ST.teamExists(ST.encodeCountryAndVal(1, 8*17)), "wrong teamExists call"
    assert not ST.teamExists(ST.encodeCountryAndVal(1, 8*17+1)), "wrong teamExists call"
    assert ST.playerExists(ST.encodeCountryAndVal(1, 8*PLAYERS_PER_TEAM_INIT)), "wrong playerExists call"
    assert ST.playerExists(ST.encodeCountryAndVal(1, 8*PLAYERS_PER_TEAM_INIT+1)), "wrong playerExists call"
    assert ST.playerExists(ST.encodeCountryAndVal(1, 8*17*PLAYERS_PER_TEAM_INIT)), "wrong playerExists call"
    assert not ST.playerExists(ST.encodeCountryAndVal(1, 8*17*PLAYERS_PER_TEAM_INIT+1)), "wrong playerExists call"

    assert ST.getTeamIdxInCountryFromPlayerIdxInCountry(1) == 1, "wrong getTeamIdx"
    assert ST.getTeamIdxInCountryFromPlayerIdxInCountry(18) == 1, "wrong getTeamIdx"
    assert ST.getTeamIdxInCountryFromPlayerIdxInCountry(19) == 2, "wrong getTeamIdx"
    (teamIdxInCountry, shirtNum) =  ST.getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(19)
    assert teamIdxInCountry == 2 and shirtNum == 0, "wrong team/shirtNum"

    assert ST.getDivisionCreationDay(1,1) == 0, "Wrong creation time"
    assert ST.getDivisionCreationDay(1,2) == 0, "Wrong creation time"

    assert ST.verseToUnixMonths(0) == DEPLOYMENT_IN_UNIX_MONTHS, "wrong verse to months"
    assert ST.verseToUnixMonths(10) == DEPLOYMENT_IN_UNIX_MONTHS, "wrong verse to months"
    assert ST.verseToUnixMonths(VERSES_PER_DAY*30) == DEPLOYMENT_IN_UNIX_MONTHS, "wrong verse to months"
    assert ST.verseToUnixMonths(VERSES_PER_DAY*31) == DEPLOYMENT_IN_UNIX_MONTHS + 1, "wrong verse to months"

    assert ST.getDisivionIdxFromTeamIdxInCountry(1) == 1, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry(8) == 1, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry(9) == 2, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry(TEAMS_PER_LEAGUE + TEAMS_PER_LEAGUE * LEAGUES_PER_DIVISON) == 2, "wrong divIdx"
    assert ST.getDisivionIdxFromTeamIdxInCountry(TEAMS_PER_LEAGUE + TEAMS_PER_LEAGUE * LEAGUES_PER_DIVISON + 1) == 3, "wrong divIdx"

    playerIdx = ST.encodeCountryAndVal(1,35)
    playerState = ST.getPlayerSkillsAtBirth(playerIdx)
    assert playerState.getPlayerIdx() == playerIdx, "wrong playerIdx set"
    assert all(playerState.getSkills() == [51,38,61,52,46]), "wrong skills set"
    assert playerState.getMonth() == 350, "wrong age"

    teamIdx = ST.encodeCountryAndVal(1,4)
    assert ST.isBotTeam(teamIdx) == True, "team not seen as bot"

    ST.acquireBot(teamIdx, ALICE)
    ST_CLIENT.acquireBot(teamIdx, ALICE)

    assert ST.isBotTeam(teamIdx) == False, "team not seen as human"

    playerIdx = ST.encodeCountryAndVal(1,3) # belongs to team1, of course
    assert ST.isPlayerTransferable(playerIdx), "country not started yet"
    assert ST.timeZoneUpdates[timeZone].updateCycleIdx == 0, "incorrect updateCycleIdx"

    teamIdx2 = ST.encodeCountryAndVal(1, 2)
    shouldFail(lambda x: ST.movePlayerToTeam(playerIdx, teamIdx2), "should not be able to transfer from or to Bot Teams")
    ST.acquireBot(teamIdx2, BOB)
    ST_CLIENT.acquireBot(teamIdx2, BOB)
    shouldFail(lambda x: ST.movePlayerToTeam(playerIdx, teamIdx2), "should not be able to transfer from or to Bot Teams")
    teamIdx1 = ST.encodeCountryAndVal(1, 1)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == FREEVERSE, "wrong owner of player"
    ST.acquireBot(teamIdx1, CAROL)
    ST_CLIENT.acquireBot(teamIdx1, CAROL)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == CAROL, "wrong owner of player"
    ST.movePlayerToTeam(playerIdx, teamIdx2)
    assert ST.getOwnerAddrFromPlayerIdx(playerIdx) == BOB, "wrong owner of player"

    # we are at verse = 0. The league starts at verse = 3
    for v in range(3):
        assert ST.currentTimeZoneToUpdate() == (TZ_NULL, TZ_NULL), "incorrect timeZone to update"
        assert ST.timeZoneUpdates[timeZone].updateCycleIdx == 0, "incorrect updateCycleIdx"
        advanceNVerses(1, ST, ST_CLIENT)
    assert ST.currentVerse == 3, "wrong verse num"
    for v in range(24):
        for sv in range(4):
            assert ST.currentTimeZoneToUpdate() == ((v+1) % 24, sv), "incorrect timeZone to update"
            advanceNVerses(1, ST, ST_CLIENT)
    assert ST.currentVerse == 99, "wrong verse num"
    assert ST.timeZoneUpdates[timeZone].updateCycleIdx == 4, "incorrect updateCycleIdx"

    playerIdx = ST.encodeCountryAndVal(1,12) # belongs to team1, of course
    assert not ST.isPlayerTransferable(playerIdx), "country busy playing"

    teamIdx1 = ST.encodeCountryAndVal(1, 6) # from first divison
    teamIdx2 = ST.encodeCountryAndVal(1, 60) # from second divison
    action00 = {"teamIdx": teamIdx1, "teamOrder": ORDER1, "tactics": TACTICS["433"]}
    action01 = {"teamIdx": teamIdx2, "teamOrder": ORDER2, "tactics": TACTICS["442"]}
    ST_CLIENT.accumulateAction(action00)
    ST_CLIENT.accumulateAction(action01)

    verseAtLastMatch = 3 + 13 * VERSES_PER_DAY + 4
    advanceNVerses(verseAtLastMatch-ST.currentVerse, ST, ST_CLIENT)
    assert ST.currentVerse == verseAtLastMatch, "error in advance verse"
    assert not ST.isPlayerTransferable(playerIdx), "player should be free, since country is settled"
    advanceNVerses(1, ST, ST_CLIENT)
    assert ST.isPlayerTransferable(playerIdx), "player should be free, since country is settled"

    if False:
        usersInitData = {
            "teamIdxs": [teamIdx1, teamIdx2],
            "teamOrders": [DEFAULT_ORDER, REVERSE_ORDER],
            "tactics": [TACTICS["442"], TACTICS["433"]]
        }

        # Create league in BC and CLIENT. The latter stores things pre-hash too,
        # and computes stuff that will be used in potential challenges.
        leagueIdx          = ST.createLeague(verseInit, verseStep, usersInitData)
        leagueIdx_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)

        assert (leagueIdx == leagueIdx_client), "leagueIdx not in sync BC vs client"
        assert ST.isLeagueIsAboutToStart(leagueIdx), "League not detected as created"

        # advance a bit before first match to change tactics
        assert ST.currentVerse == 0, "We should start with verse 0"
        advanceToBlock(ST.nextVerseBlock()-5, ST, ST_CLIENT)

        # receive the first action! Every time it arrives to the Client, it acumulates it
        action00 = {"teamIdx": teamIdx1, "teamOrder": ORDER1, "tactics": TACTICS["433"]}
        action01 = {"teamIdx": teamIdx2, "teamOrder": ORDER2, "tactics": TACTICS["442"]}
        ST_CLIENT.accumulateAction(action00)
        ST_CLIENT.accumulateAction(action01)

        # Advance to just before matchday 2, which starts at verse 3 + 24 = 27
        # From verse 0 to 26:
        assert ST.currentVerse == 0, "We should start with verse 0"
        advanceNVerses(24, ST, ST_CLIENT)
        assert ST.currentVerse == 24, "We should be at verse 24, league finishes at 27"
        advanceToBlock(ST.nextVerseBlock()-5, ST, ST_CLIENT)

        assert ST.hasLeagueStarted(leagueIdx), "League not detected as already being played"
        assert not ST.hasLeagueFinished(leagueIdx), "League not detected as not finished yet"

        # Cook data to change tactics before games in matchday 2 begin.
        action0 = {"teamIdx": teamIdx1, "teamOrder": ORDER2, "tactics": TACTICS["433"]}
        action1 = {"teamIdx": teamIdx2, "teamOrder": DEFAULT_ORDER, "tactics": TACTICS["442"]}

        ST_CLIENT.accumulateAction(action0)
        advanceNVerses(2, ST, ST_CLIENT)
        assert ST.currentVerse == 26, "We should be at verse 26, league finishes at 27"
        advanceToBlock(ST.nextVerseBlock()-5, ST, ST_CLIENT)
        ST_CLIENT.accumulateAction(action1)

        assert not ST.hasLeagueFinished(leagueIdx), "League detected as finished when it is still being played"
        assert not ST.hasLeagueBeenUpdated(leagueIdx), "League was updated too early, before finishing"

        # Move beyond league end, and force a lie by the updater:
        ST_CLIENT.forceVerseRootLie = True
        advanceToBlock(ST.nextVerseBlock()+1, ST, ST_CLIENT)
        ST_CLIENT.forceVerseRootLie = False

        assert ST.hasLeagueFinished(leagueIdx), "League not detected as already finished"
        assert ST.hasLeagueBeenUpdated(leagueIdx), "League not detected as updated, when the sync process should have done it"

        # this is the verse number that we will be challenging:
        verse = ST.leagues[leagueIdx].verseFinal()

        # First check that the status is correct
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL1) # Level 1  (lie)

        # Challenge with the truth
        challengeLevel1(verse, BOB, ST, ST_CLIENT, lie = 0)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # Level 2 (truth)

        # Try to challenge by providing the truth too... it must fail, since the Blockchain detects
        # that you're making a statement compatible with the previous update.
        # Basically, the merkle root of your new data equals the hash that you're challenging.
        superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
        subVerse = 0
        shouldFail(lambda x: challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie = 0), \
            "You were able to challenge a superroot by providing compatible data")
        # so the state remains the same:
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # Level 2 (truth)


        # Try to challenge by providing a Lie about one of the leagues root, it will be caught later on
        challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie=2)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (lie)

        # A Challenger provides... yet another a lie at matchday 0
        dayToLie = 0
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= dayToLie + 2)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (lie)

        # it is caught instantly, which sends us back to one level up
        selectedMatchday = 0
        challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (lie)

        # A Challenger provides... yet another lie at matchday 1
        dayToLie = 1
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= dayToLie + 2)

        # it is caught instantly, which sends us back to one level up
        selectedMatchday = 1
        challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (lie)

        # A Challenger provides... yet another lie at initskills
        dayToLie = -1 # -1 corresponds to initSkills
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= dayToLie + 2)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (lie)

        # it is caught instantly, which sends us back to one level up
        challengeLevel4(LEAGUE_INIT_SKILLS_ID, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (lie)

        # A Challenger finally provides the truth that proves that Level 3 was a lie
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= 0)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (truth)

        # every challenge to this update will fail instantly
        challengeLevel4(LEAGUE_INIT_SKILLS_ID, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (truth)

        selectedMatchday = 0
        challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (truth)

        selectedMatchday = 1
        challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (truth)

        # at this point we basically know that the level 4 update was TRUE (so level 3 was lying).
        # To prove it, some time passes, and the status changes
        advanceNBlocks(CHALLENGING_PERIOD_BLKS + 1, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # back to Level 2 (truth)

        # We can now see that we should lash the guy at level 3, now that time has passed.
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert needsSlash == UPDT_LEVEL3, "We should be able to slash AllLeagues, but not detected"

        # important: we could do this slash manually (now that time has passed).
        # regardless of this, it will be done automatically if now someone challenges (again)
        # the current Level 2 data. Let's see it:

        # Good, so we're now at Level 2, which was true, let's challenge again... with a lie.
        challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie=3)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (lie)

        # Check that the previous guy was already slashed (and all his update data was erased):
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert needsSlash == UPDT_NONE, "The previous challenge shouldve slashed AllLeague, but didnot"

        # Good, let catch this lie by telling the truth:
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= 0)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (truth)

        # time passes (no-one would dare to challenge this truth again)
        advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # back to Level 2 (truth)

        # check that we can manually slash the Level 3 guy:
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert needsSlash == UPDT_LEVEL3, "We should be able to slash AllLeagues, but not detected"

        # check that the verse is not yet settled, but if some times passes, it will settle.
        assert not isVerseSettled, "Verse incorrectly detected as settled"
        advanceNBlocks(CHALLENGING_PERIOD_BLKS, ST, ST_CLIENT)
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert isVerseSettled, "Verse incorrectly detected as not yet settled"

        # now that it is settled, no challenge is accepted:
        shouldFail(lambda x: ST.assertCanChallengeStatus(verse, UPDT_LEVEL2),\
                         "League is settled, but not detected")


        # OK: let's start New League To Test that SuperRoot lies can be caught too
        verseInit = ST.currentVerse + 2
        leagueIdx          = ST.createLeague(verseInit, verseStep, usersInitData)
        leagueIdx_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)
        assert (leagueIdx == leagueIdx_client), "leagueIdx not in sync BC vs client"
        assert ST.isLeagueIsAboutToStart(leagueIdx), "League not detected as created"

        # Advance to end of league and submit a truth
        assert ST_CLIENT.forceVerseRootLie == False, "we want to tell the truth in the automatic updates now!"
        advanceToEndOfLeague(leagueIdx, ST, ST_CLIENT)
        verse = ST.leagues[leagueIdx].verseFinal()
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL1) # Level 1 (truth)

        # Check that I cannot challenge a truth with another truth:
        # (again, the merkle root of my provided data coincides with the hash that I am challenging)
        superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
        shouldFail(lambda x: challengeLevel1(verse, BOB, ST, ST_CLIENT, lie = 0), "Updater should have lied in superroot, but didnt")

        # Challenge with a lie:
        challengeLevel1(verse, BOB, ST, ST_CLIENT, lie = 4)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # Level 2 (lie)

        # Challenge with truth:
        challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie=0)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (truth)

        # Check that no-one needs to be slashed until some time passes
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert needsSlash == UPDT_NONE, "Verse incorrectly reporting slash needed"
        assert not isVerseSettled, "Verse incorrectly detected as settled"

        # Since verse is not settled yet, we cannot transfer players involved in this league. Check it:
        shouldFail(lambda x:
            exchangePlayers(
                ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1),
                ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6),
                ST, ST_CLIENT
            ), "Player sell/buy was allowed but previous league was not settled yet"
         )

        # So let's wait enough...
        advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
        # ...and check that we moved up to level 1:
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL1) # Level 1 (truth)
        # ...and that we should manually slash the level 2 lier:
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert needsSlash == UPDT_LEVEL2, "Verse incorrectly reporting slash not needed"

        # The verse is not settled yet! We can still challenge if we wanted.
        assert not isVerseSettled, "Verse incorrectly detected as already settled"
        # ...until some more time passes:
        advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert isVerseSettled, "Verse incorrectly detected as not yet settled"

        # Skip this part if you are only interested in Stakers games
        # Now, finally, we can transfer players.
        # Let's check that their teams are correctly reported before and after
        playerIdx1 = ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1)
        playerIdx2 = ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6)
        # we can either exchange players:
        exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT)
        # or transfer them one by one:
        playerIdx3 = ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx3, 2)
        team3, shirt3 = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx3)
        assert team3 == teamIdx3, "some is wrong with team assignments"
        movePlayerToTeam(playerIdx3, teamIdx1, ST, ST_CLIENT)
        team, shirt = ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx3)
        assert team == teamIdx1, "wrong initial assignment"
        assert ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx1) == (teamIdx4, 24), "Exchange did not register properly in BC"
        assert ST.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx2) == (teamIdx1, 24), "Exchange did not register properly in BC"
        assert ST_CLIENT.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx1) == (teamIdx4, 24), "Exchange did not register properly in BC"
        assert ST_CLIENT.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx2) == (teamIdx1, 24), "Exchange did not register properly in BC"


        #           -----  LEAGUE 3 ------
        # After the player exchange... we create another league
        verseInit = ST.currentVerse + 2
        leagueIdx          = ST.createLeague(verseInit, verseStep, usersInitData)
        leagueIdx_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)
        assert leagueIdx == leagueIdx_client, "Leagues in client not in sync with BC"

        # Move to end of league, during the challenging period. The Synchronizer has updated the verse with truth.
        advanceToEndOfLeague(leagueIdx, ST, ST_CLIENT)
        verse = ST_CLIENT.leagues[leagueIdx].verseFinal()
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL1) # Level 1 (truth)
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert not isVerseSettled, "Verse incorrectly detected as settled"

        # Try to challenge by providing a false ALL-LEAGUES
        challengeLevel1(verse, BOB, ST, ST_CLIENT, lie = 4)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # Level 2 (lie)

        # and yet another challenge with lie:
        challengeLevel2(verse, subVerse, BOB, ST, ST_CLIENT, lie = 5)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (lie)

        # Try to challenge by providing a false ONE-LEAGUE...
        dayToLie = -1 # this is a lie in initSkills
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= dayToLie + 2)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (lie)

        # We successfully challenge the ONE-LEAGUE, and return to ALL-LEAGUES
        challengeLevel4(LEAGUE_INIT_SKILLS_ID, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3)  # Level 3 (lie)

        # We now successfully challenge the false ALL-LEAGUES by telling the truth
        challengeLevel3(verse, leagueIdx, CAROL, ST, ST_CLIENT, lie= 0)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (true)

        # We fail to prove that anything was wrong
        challengeLevel4(LEAGUE_INIT_SKILLS_ID, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (true)

        # it also fails at proving that any matchday is wrong
        selectedMatchday = 0
        challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (true)
        selectedMatchday = 1
        challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL4) # Level 4 (true)

        # finally, after a CHLL_PERIOD, it shows that it is back to the superRoot
        advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL2) # Level 2 (lie)
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert not isVerseSettled, "Verse incorrectly detected as settled"

        # we challenge with truth:
        challengeLevel2(verse, subVerse, CAROL, ST, ST_CLIENT, lie = 0)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL3) # Level 3 (truth)

        # we wait and we see that we're back to Level 1, and should slash the level 2 lier
        advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
        ST.assertCanChallengeStatus(verse, UPDT_LEVEL1) # Level 1 (truth)
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert needsSlash == UPDT_LEVEL2, "we did not detect that this guy should be slashed"
        assert not isVerseSettled, "Verse incorrectly detected as settled"

        # the verse settles after time:
        advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert isVerseSettled, "Verse incorrectly detected as not yet settled"
        assert needsSlash == UPDT_LEVEL2, "we did not detect that this guy should be slashed"


        # We make sure that we can inquire the state of any player after these leagues and player sales:
        assertPlayerStateInClientIsCertifiable(1, ST, ST_CLIENT)

        # The following all-team printout is interesting. On the one hand, it checks that all player states
        # in that team can be certified by the BC. On the other hand, you can check that the 2nd player
        # corresponds to the player bought from team4, in the exchange done above. Basically, that transfers went right.
        printTeam(teamIdx1, ST_CLIENT)

        # NOW: GO WILD
        # create many teams, and leagues, and mess it all.
        advanceNVerses(1000, ST, ST_CLIENT)
        nTeams      = 200
        nLeagues    = 20
        nPlayers    = 400

        # create many teams
        for t in range(nTeams):
            createTeam("BotTeam"+str(t), ALICE, ST, ST_CLIENT)

        # transfer many players
        for p in range(nPlayers):
            playerIdx1 = 1+intHash(str(p)) % 100*PLAYERS_PER_TEAM_MAX
            playerIdx2 = 1+intHash(str(p)+ "salt") % 100 * PLAYERS_PER_TEAM_MAX
            exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT)
            pylio.assertPlayerStateInClientIsCertifiable(playerIdx1, ST, ST_CLIENT)

        lastTeamIdx = 1
        nTeamsPerLeague = 8

        # Create Lots of Leagues for testing:
        leaguesTested = []
        for l in range(nLeagues):
            verseInit = ST.currentVerse + 7
            usersInitData = {
                "teamIdxs": [t for t in range(lastTeamIdx,lastTeamIdx+nTeamsPerLeague)],
                "teamOrders": [getRandomElement(POSSIBLE_ORDERS, t) for t in range(lastTeamIdx,lastTeamIdx+nTeamsPerLeague)],
                "tactics": [getRandomElement(POSSIBLE_TACTICS, t) for t in range(lastTeamIdx,lastTeamIdx+nTeamsPerLeague)]
            }
            lastTeamIdx += nTeamsPerLeague
            leagueIdx = ST.createLeague(verseInit, verseStep, usersInitData)
            leagueIdx_client = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)
            leaguesTested.append(leagueIdx)

            assert (leagueIdx == leagueIdx_client), "leagueIdx not in sync BC vs client"
            assert ST.isLeagueIsAboutToStart(leagueIdx), "League not detected as created"
            advanceNVerses(intHash(str(l))%2 , ST, ST_CLIENT) # advance either 1 or 0 verses


        # We run a large set of updates, challenges, exchanges, etc.
        # Only the verseRoot is true in all of them.
        # After the "brutalBlock" call, we may end up in one of the many possible verse states, depending
        # on the lying sequence.
        # But after the brutalBlock, we will still be in the challenging period... and hence, after telling
        # the truths in all such leagues, we should end up with all of them in the original Level 1 honest state.
        ST, ST_CLIENT = brutalBlock(ST, ST_CLIENT, leaguesTested)
        # We now tell the truth for all leagues. We need to do this in 2 steps, since some leagues
        # ended up level 3 or 4. That's why we wait for an extra challenging block, and repeat telling truth.
        # By the way, we also try to transfer many players at the same time by setting the "True" bool
        ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, True)
        advanceNBlocks(CHALLENGING_PERIOD_BLKS + 1, ST, ST_CLIENT)
        ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, True)

        # Wait for everything to settle and check we're in the SuperRoot state in all such leagues
        advanceNVerses(2000, ST, ST_CLIENT)
        for leagueIdx in leaguesTested:
            verse = ST.leagues[leagueIdx].verseFinal()
            verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
            assert isVerseSettled
            assert verseStatus == UPDT_LEVEL1, "league should be back to Verse..."

        # Repeat the brutal sequence.
        ST, ST_CLIENT = brutalBlock(ST, ST_CLIENT, leaguesTested)
        ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, False)
        advanceNBlocks(CHALLENGING_PERIOD_BLKS + 1, ST, ST_CLIENT)
        ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, True)
        advanceNVerses(2000, ST, ST_CLIENT)
        for leagueIdx in leaguesTested:
            verse = ST.leagues[leagueIdx].verseFinal()
            verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
            assert isVerseSettled
            assert verseStatus == UPDT_LEVEL1, "league should be back to Verse..."


        # If we made it to this point, the tests basically passed.
        # We will return the hash of the entire ST and ST_CLIENT gigantic structs.
        # This will ensure that not only is the flow as expected, but the stored data
        # does not change accidentally.

    testResult = intHash(serialize2str(ST) + serialize2str(ST_CLIENT)) % 1000
    return testResult


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


