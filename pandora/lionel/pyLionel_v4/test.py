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

def updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, doExchanges):
    # Now go through every league and verse, and always tell the truth in all challenges.
    # Since the verseRoot was true in all cases, we should end up in verseRoot state at the very end.
    for extraVerse in range(2000):
        if doExchanges and extraVerse % 10:
            for p in range(2):
                playerIdx1 = 1 + intHash(str(p+extraVerse)) % 100 * NPLAYERS_PER_TEAM
                playerIdx2 = 1 + intHash(str(p+extraVerse) + "salt") % 100 * NPLAYERS_PER_TEAM
                try:
                    exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT)
                except:
                    pass
                pylio.assertPlayerStateInClientIsCertifiable(playerIdx1, ST, ST_CLIENT)

        for leagueIdx in leaguesTested:
            verse = ST.leagues[leagueIdx].verseFinal()
            verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
            assert not ( ST.isLeagueSettled(leagueIdx) and verseStatus != UPDT_VERSE), "Someone hacked the game"
            if ST.hasLeagueFinished(leagueIdx) and (not ST.isLeagueSettled(leagueIdx)):
                print("challenging league...", leagueIdx)
                if verseStatus == UPDT_SUPROOTS:
                    print("challenging league... superRoot", leagueIdx)
                    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
                    subVerse = 0
                    ST.challengeSuperRoot(verse, subVerse, leagueRoots[subVerse], ADDR2)
                    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)
                if verseStatus == UPDT_LGROOTS:
                    print("challenging league... leagueRoot", leagueIdx)
                    dataToChallengeLeague = ST_CLIENT.leagues[leagueIdx].dataToChallengeLeague
                    ST.challengeleagueRoots(
                        verse,
                        ST.getPosInSubverse(verse, leagueIdx),
                        dataToChallengeLeague,
                        ADDR3
                    )
                elif verseStatus == UPDT_ONELEAGUE:
                    thisLeagueIdx = ST.getLeagueIdxFromPosInSubverse(verse, ST.verseToLeagueCommits[verse].posInSubVerse)
                    print("challenging league... initSkills", thisLeagueIdx)
                    ST.challengeInitSkills(
                        verse,
                        ST_CLIENT.leagues[thisLeagueIdx].usersInitData,
                        duplicate(ST_CLIENT.leagues[thisLeagueIdx].dataToChallengeInitSkills)
                    )
        # advanceNVerses(2, ST, ST_CLIENT)
        return ST, ST_CLIENT

def brutalBlock(ST, ST_CLIENT, leaguesTested):
    # We will loop for leagues so that we first always lie:
    #   - first: challenge the superRoot with lie
    #   - second: challenge the allLeagues with lie
    # Then catch the last lie with a challengeInitSkills
    # And then always tell the truth
    # We keep a list of leagues that already went down to the challengeInitskills to know
    # that we now need to tell the truth.
    leaguesTestedAtLevel3 = []
    advanceNVerses(250, ST, ST_CLIENT)
    for extraVerse in range(45):
        advanceNVerses(1, ST, ST_CLIENT)
        for leagueIdx in leaguesTested:
            if ST.hasLeagueFinished(leagueIdx) and (not ST.isLeagueSettled(leagueIdx)):
                print("challenging league...", leagueIdx)
                verse = ST.leagues[leagueIdx].verseFinal()
                verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
                if verseStatus == UPDT_VERSE:
                    print("challenging league... verseRoot", leagueIdx)
                    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
                    superRootsLie, leagueRootsLie = pylio.createLieSuperRoot(superRoots, leagueRoots, 11)
                    ST.challengeVerseRoot(verse, superRootsLie, ADDR2)
                    ST.assertCanChallengeStatus(verse, UPDT_SUPROOTS)
                elif verseStatus == UPDT_SUPROOTS:
                    print("challenging league... superRoot", leagueIdx)
                    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
                    superRootsLie, leagueRootsLie = pylio.createLieSuperRoot(superRoots, leagueRoots, 12)
                    subVerse = 0
                    ST.challengeSuperRoot(verse, subVerse, leagueRootsLie[subVerse], ADDR2)
                    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)
                elif verseStatus == UPDT_LGROOTS:
                    if leagueIdx in leaguesTestedAtLevel3:
                        print("challenging league... allLeagues with truth: ", leagueIdx)
                        dataToChallengeLeague = ST_CLIENT.leagues[leagueIdx].dataToChallengeLeague
                        ST.challengeleagueRoots(
                            verse,
                            ST.getPosInSubverse(verse, leagueIdx),
                            dataToChallengeLeague,
                            ADDR3
                        )
                        ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
                    else:
                        print("challenging league... allLeagues with lie: ", leagueIdx)
                        dataToChallengeLeague = ST_CLIENT.leagues[leagueIdx].dataToChallengeLeague
                        dataToChallengeLeagueLie = pylio.duplicate(dataToChallengeLeague)
                        dataToChallengeLeagueLie.initSkillsHash += 1
                        dataToChallengeLeagueLie.dataAtMatchdayHashes[0] += 1
                        ST.challengeleagueRoots(
                            verse,
                            ST.getPosInSubverse(verse, leagueIdx),
                            dataToChallengeLeagueLie,
                            ADDR3
                        )
                        ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

                elif verseStatus == UPDT_ONELEAGUE:
                    thisLeagueIdx = ST.getLeagueIdxFromPosInSubverse(verse, ST.verseToLeagueCommits[verse].posInSubVerse)
                    print("challenging league... initSkills", thisLeagueIdx)
                    ST.challengeInitSkills(
                        verse,
                        ST_CLIENT.leagues[thisLeagueIdx].usersInitData,
                        duplicate(ST_CLIENT.leagues[thisLeagueIdx].dataToChallengeInitSkills)
                    )
                    if thisLeagueIdx in leaguesTestedAtLevel3:
                        ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
                    else:
                        ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)
                    leaguesTestedAtLevel3.append(thisLeagueIdx)
    return ST, ST_CLIENT



# TEST1: create a team, print players
# Exchange 2 players in different teams, check that all is updated OK
# the test is passed if the hash mod 1000 of all that is printed is as expected
def test1():
    ST          = Storage(isClient = False)
    ST_CLIENT   = Storage(isClient = True)

    teamIdx1 = ST.createTeam("Barca", ADDR1)
    teamIdx2 = ST.createTeam("Madrid", ADDR2)

    teamIdx1_client = ST_CLIENT.createTeam("Barca", ADDR1)
    teamIdx2_client = ST_CLIENT.createTeam("Madrid", ADDR2)

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "TeamIdx not in sync BC vs client"

    # Test that we can ask the BC if state of a player (computed by the Client) is correct:
    pylio.assertPlayerStateInClientIsCertifiable(1, ST, ST_CLIENT)

    print("Team created with teamIdx, teamName = " + str(teamIdx1) + ", " + ST.teams[teamIdx1].name)
    hash0 = printTeam(teamIdx1, ST_CLIENT)

    print("\n\nplayers 2 and 24 before sale:\n")

    hash1 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(2))

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "PlayerStates not in sync BC vs client"

    print("\n")
    hash2 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(24))

    advanceNBlocks(10, ST, ST_CLIENT)

    exchangePlayers(2, 24, ST, ST_CLIENT)

    print("\n\nplayers 2 and 24 after sale:\n")
    hash3 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(2))
    print("\n")
    hash4 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(24))
    hashSum         = hash0+hash1+hash2+hash3+hash4
    return hashSum



# TEST2: all the workflow!
def test2():
    # Create contract storage in BC, and its extended version for the CLIENT
    # We need to keep both in sync. The CLIENT stores, besides what is in the BC, the pre-hash stuff.
    ST          = Storage(isClient = False)
    ST_CLIENT   = Storage(isClient = True)

    # The accumulator is responsible for receiving user actions and committing them in the correct verse.
    ST_CLIENT.addAccumulator()

    # Note that every 'advance' we do will check if some user actions need to be commited
    advanceToBlock(10, ST, ST_CLIENT)

    # Create teams in BC and client
    teamIdx1 = ST.createTeam("Barca", ADDR1)
    teamIdx2 = ST.createTeam("Madrid", ADDR2)
    teamIdx3 = ST.createTeam("Milan", ADDR3)
    teamIdx4 = ST.createTeam("PSG", ADDR3)

    teamIdx1_client = ST_CLIENT.createTeam("Barca", ADDR1)
    teamIdx2_client = ST_CLIENT.createTeam("Madrid", ADDR2)
    teamIdx3_client = ST_CLIENT.createTeam("Milan", ADDR3)
    teamIdx4_client = ST_CLIENT.createTeam("PSG", ADDR3)

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "TeamIdx not in sync BC vs client"

    # advances both BC and CLIENT, and syncs if it goes through a verse
    advanceToBlock(100, ST, ST_CLIENT)

    # One verse is about 1 hour, so a day is about 24 verseSteps
    verseInit = 3
    verseStep = 24

    # Cook init data for the 1st league
    usersInitData = {
        "teamIdxs": [teamIdx1, teamIdx2],
        "teamOrders": [DEFAULT_ORDER, REVERSE_ORDER],
        "tactics": [TACTICS["442"], TACTICS["433"]]
    }

    # Create league in BC and CLIENT. The latter stores things pre-hash too
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

    verse = ST.leagues[leagueIdx].verseFinal()

    # Since the entire verse was updated faithfully, any challenge to it will fail.
    # First check that the status is correct
    ST.assertCanChallengeStatus(verse, UPDT_VERSE)

    # Challenge with a lie
    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
    ST.challengeVerseRoot(verse, superRoots, ADDR2)


    # Try to challenge by providing a correct All-leagues-roots... should fail
    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
    ST_CLIENT.getSubVerseData(verse)
    subVerse = 0
    pylio.shouldFail(lambda x: ST.challengeSuperRoot(verse, subVerse, leagueRoots[subVerse], ADDR2), \
        "You were able to challenge a superroot by providing compatible leagueRoots")

    # Try to challenge by providing a lie about one of the leagues root, it will be caught later on
    superRootsLie, leagueRootsLie = pylio.createLieSuperRoot(superRoots, leagueRoots, 2)
    ST.challengeSuperRoot(verse, subVerse, leagueRootsLie[subVerse], ADDR2)

    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)

    # A Challenger provides a lie at matchday 0
    dataToChallengeLeague = ST_CLIENT.leagues[leagueIdx].dataToChallengeLeague
    dataToChallengeLeagueLie = pylio.duplicate(dataToChallengeLeague)
    dataToChallengeLeagueLie.dataAtMatchdayHashes[0] += 1

    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeagueLie,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    selectedMatchday = 0
    challengeLeagueAtSelectedMatchday(selectedMatchday, verse, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)

    # A Challenger provides a lie at matchday 1
    dataToChallengeLeagueLie = pylio.duplicate(dataToChallengeLeague)
    dataToChallengeLeagueLie.dataAtMatchdayHashes[1] += 1

    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeagueLie,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    selectedMatchday = 1
    challengeLeagueAtSelectedMatchday(selectedMatchday, verse, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)


    # A Challenger provides a lie at initskills
    dataToChallengeLeagueLie = pylio.duplicate(dataToChallengeLeague)
    dataToChallengeLeagueLie.initSkillsHash += 1

    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeagueLie,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    ST.challengeInitSkills(
        verse,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
    )
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)

    # A Challenger provides the truth
    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeague,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

    # every further challenge fails
    ST.challengeInitSkills(
        verse,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    selectedMatchday = 0
    challengeLeagueAtSelectedMatchday(selectedMatchday, verse, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    selectedMatchday = 1
    challengeLeagueAtSelectedMatchday(selectedMatchday, verse, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

    # at this point we basically know that the provided Matchdays data is wrong.
    # to prove it, some time passes, and the status changes
    advanceNBlocks(CHALLENGING_PERIOD_BLKS + 1, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_SUPROOTS)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert needsSlash == UPDT_LGROOTS, "We should be able to slash AllLeagues, but not detected"

    superRootsLie, leagueRootsLie = pylio.createLieSuperRoot(superRoots, leagueRoots, 3)
    ST.challengeSuperRoot(verse, subVerse, leagueRootsLie[subVerse], ADDR2)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert needsSlash == UPDT_NONE, "The previous challenge shouldve slashed AllLeague, but didnot"
    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeague,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert needsSlash == UPDT_LGROOTS, "We should be able to slash AllLeagues, but not detected"
    assert not isVerseSettled, "Verse incorrectly detected as settled"
    ST.assertCanChallengeStatus(verse, UPDT_SUPROOTS)
    advanceNBlocks(CHALLENGING_PERIOD_BLKS, ST, ST_CLIENT)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert isVerseSettled, "Verse incorrectly detected as not yet settled"
    pylio.shouldFail(lambda x: ST.assertCanChallengeStatus(verse, UPDT_SUPROOTS),\
                     "League is settled, but not detected")


    # Create a New League To Test that SuperRoot lies can be caught too
    verseInit = ST.currentVerse + 2
    leagueIdx          = ST.createLeague(verseInit, verseStep, usersInitData)
    leagueIdx_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)
    assert (leagueIdx == leagueIdx_client), "leagueIdx not in sync BC vs client"
    assert ST.isLeagueIsAboutToStart(leagueIdx), "League not detected as created"

    # Advance to end of league and submit a truth
    verse = ST.leagues[leagueIdx].verseFinal()
    advanceToEndOfLeague(leagueIdx, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_VERSE)

    # Check that a lie can be caught by comparing with local computation
    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
    pylio.shouldFail(lambda x: ST.challengeVerseRoot(verse, superRoots, ADDR2), "Updater should have lied in superroot, but didnt")
    superRootsLie, leagueRootsLie = createLieSuperRoot(superRoots, leagueRoots, 4)

    # Submit a challenge and check its time evolution after waiting....
    ST.challengeVerseRoot(verse, superRootsLie, ADDR2)
    ST.assertCanChallengeStatus(verse, UPDT_SUPROOTS)
    ST.challengeSuperRoot(verse, subVerse, leagueRoots[subVerse], ADDR2)
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert needsSlash == UPDT_NONE, "Verse incorrectly reporting slash needed"
    assert not isVerseSettled, "Verse incorrectly detected as settled"

    # We do not wait enough and try a sell/buy action is attempted
    pylio.shouldFail(lambda x:
        ST.exchangePlayers(
            ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1), ADDR1,
            ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6), ADDR3
        ), "Player sell/buy was allowed but previous league was not settled yet"
     )

    # We now wait enough
    advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    ST.assertCanChallengeStatus(verse, UPDT_VERSE)
    assert needsSlash == UPDT_SUPROOTS, "Verse incorrectly reporting slash not needed"
    assert not isVerseSettled, "Verse incorrectly detected as already settled"
    advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert isVerseSettled, "Verse incorrectly detected as not yet settled"

    # Now, we can sell-buy players
    playerIdx1 = ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1)
    playerIdx2 = ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6)

    ST.exchangePlayers(
        playerIdx1, ADDR1,
        playerIdx2, ADDR3
    )
    ST_CLIENT.exchangePlayers(
        playerIdx1, ADDR1,
        playerIdx2, ADDR3
    )

    assert ST.getTeamIdxAndShirtForPlayerIdx(playerIdx1) == (teamIdx4, 6), "Exchange did not register properly in BC"
    assert ST.getTeamIdxAndShirtForPlayerIdx(playerIdx2) == (teamIdx1, 1), "Exchange did not register properly in BC"
    assert ST_CLIENT.getTeamIdxAndShirtForPlayerIdx(playerIdx1) == (teamIdx4, 6), "Exchange did not register properly in BC"
    assert ST_CLIENT.getTeamIdxAndShirtForPlayerIdx(playerIdx2) == (teamIdx1, 1), "Exchange did not register properly in BC"

    #           -----  LEAGUE 2 ------
    # After the player exchange... we create another league
    verseInit = ST.currentVerse + 2
    leagueIdx          = ST.createLeague(verseInit, verseStep, usersInitData)
    leagueIdx_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)
    assert leagueIdx == leagueIdx_client, "Leagues in client not in sync with BC"

    # At the end of league until, to the challenging period
    advanceToEndOfLeague(leagueIdx, ST, ST_CLIENT)
    verse = ST_CLIENT.leagues[leagueIdx].verseFinal()
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert not isVerseSettled, "Verse incorrectly detected as settled"

    # Try to challenge by providing a false ALL-LEAGUES
    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
    superRootsLie, leagueRootsLie = pylio.createLieSuperRoot(superRoots, leagueRoots, 4)
    ST.challengeVerseRoot(verse, superRootsLie, ADDR2)
    superRootsLie, leagueRootsLie = pylio.createLieSuperRoot(superRoots, leagueRoots, 5)
    ST.challengeSuperRoot(verse, subVerse, leagueRootsLie[subVerse], ADDR2)
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)

    # Try to challenge by providing a false ONE-LEAGUE...
    dataToChallengeLeague = ST_CLIENT.leagues[leagueIdx].dataToChallengeLeague
    dataToChallengeLeagueLie = pylio.duplicate(dataToChallengeLeague)
    dataToChallengeLeagueLie.initSkillsHash += 1

    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeagueLie,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

    # We successfully challenge the ONE-LEAGUE, and return to ALL-LEAGUES
    ST.challengeInitSkills(
        verse,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
    )
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)

    # We now successfully challenge the false ALL-LEAGUES
    ST.challengeleagueRoots(
        verse,
        ST.getPosInSubverse(verse, leagueIdx),
        dataToChallengeLeague,
        ADDR3
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

    # We fail to prove that anything was wrong
    ST.challengeInitSkills(
        verse,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
    )
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

    # it also fails at proving that any matchday is wrong
    selectedMatchday = 0
    challengeLeagueAtSelectedMatchday(selectedMatchday, verse, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)
    selectedMatchday = 1
    challengeLeagueAtSelectedMatchday(selectedMatchday, verse, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_ONELEAGUE)

    # finally, after a CHLL_PERIOD, it shows that it is back to the superRoot
    advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_SUPROOTS)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert not isVerseSettled, "Verse incorrectly detected as settled"

    ST.challengeSuperRoot(verse, subVerse, leagueRoots, ADDR3)
    ST.assertCanChallengeStatus(verse, UPDT_LGROOTS)
    advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
    ST.assertCanChallengeStatus(verse, UPDT_VERSE)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert not isVerseSettled, "Verse incorrectly detected as settled"
    advanceNBlocks(CHALLENGING_PERIOD_BLKS+1, ST, ST_CLIENT)
    verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
    assert isVerseSettled, "Verse incorrectly detected as not yet settled"


    # We make sure that we can inquire the state of any player after these leagues and player sales:
    pylio.assertPlayerStateInClientIsCertifiable(1, ST, ST_CLIENT)

    # The following all-team printout is interesting. On the one hand, it checks that all player states
    # in that team can be certified by the BC. On the other hand, you can check that the 2nd player
    # corresponds to the player bought from team4, in the exchange done above.
    printTeam(teamIdx1, ST_CLIENT)


    # create many teams, and leagues, and mess it all.
    advanceNVerses(1000, ST, ST_CLIENT)
    nTeams      = 200
    nLeagues    = 20
    nPlayers    = 400

    for t in range(nTeams):
        ST.createTeam("BotTeam"+str(t), ADDR1)
        ST_CLIENT.createTeam("BotTeam"+str(t), ADDR1)

    for p in range(nPlayers):
        playerIdx1 = 1+intHash(str(p)) % 100*NPLAYERS_PER_TEAM
        playerIdx2 = 1+intHash(str(p)+ "salt") % 100 * NPLAYERS_PER_TEAM
        ST.exchangePlayers(
            playerIdx1, ST.getOwnerAddrFromPlayerIdx(playerIdx1),
            playerIdx2, ST.getOwnerAddrFromPlayerIdx(playerIdx2)
        )
        ST_CLIENT.exchangePlayers(
            playerIdx1, ST_CLIENT.getOwnerAddrFromPlayerIdx(playerIdx1),
            playerIdx2, ST_CLIENT.getOwnerAddrFromPlayerIdx(playerIdx2)
        )
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
    # Only the verseRoot is true. So we may end up in one of the many other states, but false.
    ST, ST_CLIENT = brutalBlock(ST, ST_CLIENT, leaguesTested)
    # We now tell the truth for all leagues. We need to do this in 2 steps, since some leagues
    # ended up in the last state. That's why we wait for an extra challenging block, and repeat telling truth.
    ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, True)
    advanceNBlocks(CHALLENGING_PERIOD_BLKS + 1, ST, ST_CLIENT)
    ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, True)

    # Wait for everything to settle and check we're in the SuperRoot state in all such leagues
    advanceNVerses(2000, ST, ST_CLIENT)
    for leagueIdx in leaguesTested:
        verse = ST.leagues[leagueIdx].verseFinal()
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert isVerseSettled
        assert verseStatus == UPDT_VERSE, "league should be back to Verse..."


    ST, ST_CLIENT = brutalBlock(ST, ST_CLIENT, leaguesTested)
    ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, False)
    advanceNBlocks(CHALLENGING_PERIOD_BLKS + 1, ST, ST_CLIENT)
    ST, ST_CLIENT = updateAllLeaguesWithTruth(ST, ST_CLIENT, leaguesTested, True)

    # Wait for everything to settle and check we're in the SuperRoot state in all such leagues
    advanceNVerses(2000, ST, ST_CLIENT)
    for leagueIdx in leaguesTested:
        verse = ST.leagues[leagueIdx].verseFinal()
        verseStatus, isVerseSettled, needsSlash = ST.getVerseUpdateStatus(verse)
        assert isVerseSettled
        assert verseStatus == UPDT_VERSE, "league should be back to Verse..."




    # Returns test result, to later check against expected
    testResult = intHash(serialize2str(ST) + serialize2str(ST_CLIENT)) % 1000
    return testResult


def test4():
    leafs = [1,2,3,4,5,6,"rew"]
    tree, depth = make_tree(duplicate(leafs), serialHash)
    assert depth == get_depth(tree), "Depth not computed correctly"

    # We show that this library can prove 1 leaf at a time, or (below), many
    idxsToProve = [1]
    neededHashes, values = prepareProofForIdxs(idxsToProve, tree, leafs)
    print("To prove these %i leafs you need %i hashes, in a tree with %i leafs, and depth %i" \
          % (len(idxsToProve), len(neededHashes), len(leafs), depth)
          )
    success1 = verify(root(tree), depth, values, neededHashes, serialHash, debug_print=False)

    idxsToProve = [1,2,3]
    neededHashes, values = prepareProofForIdxs(idxsToProve, tree, leafs)
    print("To prove these %i leafs you need %i hashes, in a tree with %i leafs, and depth %i" \
          % (len(idxsToProve), len(neededHashes), len(leafs), depth)
          )
    success2 = verify(root(tree), depth, values, neededHashes, serialHash, debug_print=False)

    # it is also valid in the case of a single element, where the 'neededHashes' is empty,
    # as we just need the root(tree), which is passed too
    leafs = ["prew"]
    tree, depth = make_tree(duplicate(leafs), serialHash)
    idxsToProve = [0]
    neededHashes, values = prepareProofForIdxs(idxsToProve, tree, leafs)
    assert not neededHashes, "No Hash should be needed, but you have a non empty array"
    print("To prove these %i leafs you need %i hashes, in a tree with %i leafs, and depth %i" \
          % (len(idxsToProve), len(neededHashes), len(leafs), depth)
          )
    success3 = verify(root(tree), depth, values, neededHashes, serialHash, debug_print=False)


    return success1 and success2 and success3


def runTest(name, result, expected):
    success = (result == expected)
    if success:
        print(name + ": PASSED")
    else:
        print(name + ": FAILED.  Result(%s) vs Expected(%s) " % (result, expected))
    return success


success = True
success = success and runTest(name = "Test Simple Team Creation", result = test1(), expected = 10754)
success = success and runTest(name = "Test Entire Workflow",      result = test2(), expected = 260)
# success = success and runTest(name = "Test Merkle",      result = test4(), expected = True)
if success:
    print("ALL TESTS:  -- PASSED --")
else:
    print("At least one test FAILED")


# TODO:
# remove ugly:         if type(dataToChallengePlayerState) == type(DataAtMatchday(0, 0, 0)):
# add tests for getOwner

# TODO: - less important -
# do not store scores but the hash or merkle root
# unify all iniHash, serialHash, etc
# remove need for the last matchdayHash, as we just need to test the states.