import random
import numpy as np
from copy import deepcopy as duplicate
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists
import sha3
from pickle import dumps as serialize

from constants import *
from pylio import *
from structs import *

import __builtin__ as builtin

# TEST1: create a team, print players
# Exchange 2 players in different teams, check that all is updated OK
# the test is passed if the hash mod 1000 of all that is printed is as expected
def test1():
    ST          = Storage()
    ST_CLIENT   = Storage()

    teamIdx1 = createTeam("Barca", ADDR1, ST)
    teamIdx2 = createTeam("Madrid", ADDR2, ST)

    teamIdx1_client = createTeam("Barca", ADDR1, ST_CLIENT)
    teamIdx2_client = createTeam("Madrid", ADDR2, ST_CLIENT)

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "TeamIdx not in sync BC vs client"

    # Test that we can ask the BC if state of a player (computed by the Client) is correct:
    player1State            = getLastWrittenPlayerStateFromPlayerIdx(1, ST_CLIENT)
    player1ChallengeData    = computeDataToChallengePlayerIdx(1, ST_CLIENT)
    assert isCorrectStateForPlayerIdx(player1State, player1ChallengeData, ST), "Computed player state by CLIENT is not recognized by BC.."

    print "Team created with teamIdx, teamName = " + str(teamIdx1) + ", " + ST.teams[teamIdx1].name
    hash0 = printTeam(teamIdx1, ST_CLIENT)

    print "\n\nplayers 2 and 24 before sale:\n"
    hash1 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(2, ST_CLIENT))

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "PlayerStates not in sync BC vs client"

    print "\n"
    hash2 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(24, ST_CLIENT))

    ST.advanceNBlocks(10)
    ST_CLIENT.advanceNBlocks(10)

    exchangePlayers(
        2, ADDR1,
        24, ADDR2,
        ST
    )
    exchangePlayers(
        2, ADDR1,
        24, ADDR2,
        ST_CLIENT
    )

    print "\n\nplayers 2 and 24 after sale:\n"
    hash3 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(2, ST))
    print "\n"
    hash4 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(24, ST_CLIENT))
    hashSum         = hash0+hash1+hash2+hash3+hash4
    return hashSum



# TEST2: all the workflow!
def test2():
    # Create contract storage in BC, and its extended version for the CLIENT
    # We need to keep both in sync. The CLIENT stores, besides what is in the BC, the pre-hash stuff.
    ST          = Storage()
    ST_CLIENT   = Storage()

    ST.advanceToBlock(10)
    ST_CLIENT.advanceToBlock(10)

    # Create teams in BC and client
    teamIdx1 = createTeam("Barca", ADDR1, ST)
    teamIdx2 = createTeam("Madrid", ADDR2, ST)
    teamIdx3 = createTeam("Milan", ADDR3, ST)
    teamIdx4 = createTeam("PSG", ADDR3, ST)

    teamIdx1_client = createTeam("Barca", ADDR1, ST_CLIENT)
    teamIdx2_client = createTeam("Madrid", ADDR2, ST_CLIENT)
    teamIdx3_client = createTeam("Milan", ADDR3, ST_CLIENT)
    teamIdx4_client = createTeam("PSG", ADDR3, ST_CLIENT)

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "TeamIdx not in sync BC vs client"

    # Cook init data for the 1st league
    ST.advanceToBlock(100)
    ST_CLIENT.advanceToBlock(100)

    blockInit = 190
    blockStep = 10

    usersInitData = {
        "teamIdxs": [teamIdx1, teamIdx2],
        "teamOrders": [DEFAULT_ORDER, REVERSE_ORDER],
        "tactics": [TACTICS["442"], TACTICS["433"]]
    }

    # Create league in BC and CLIENT. The latter stores things pre-hash too
    leagueIdx          = createLeague(blockInit, blockStep, usersInitData, ST)
    leagueIdx_client   = createLeagueClient(blockInit, blockStep, usersInitData, ST_CLIENT)

    assert ST.leagues[leagueIdx].isLeagueIsAboutToStart(ST.currentBlock), "League not detected as created"

    # Advance to matchday 2
    ST.advanceToBlock(blockInit + blockStep - 5)
    ST_CLIENT.advanceToBlock(blockInit + blockStep - 5)
    assert ST.leagues[leagueIdx].hasLeagueStarted(ST.currentBlock), "League not detected as already being played"
    assert not ST.leagues[leagueIdx].hasLeagueFinished(ST.currentBlock), "League not detected as not finished yet"

    # Cook data to change tactics before games in matchday 2 begin.
    # Note that we could specify only for 1 of the teams if we wanted.
    # TODO: force that this happens atomically, 1 team at each transaction
    usersAlongData = {
        "teamIdxsWithinLeague": [1],
        "tactics": [TACTICS["433"], TACTICS["442"]],
        "block": ST.currentBlock
    }


    # Submit data to change tactics
    ST.leagues[leagueIdx].updateUsersAlongDataHash(usersAlongData, ST.currentBlock)
    ST_CLIENT.leagues[leagueIdx].updateUsersAlongData(usersAlongData, ST.currentBlock)

    # Move beyond league end
    ST.advanceNBlocks(blockStep)
    ST_CLIENT.advanceNBlocks(blockStep)
    assert ST.leagues[leagueIdx].hasLeagueFinished(ST.currentBlock), "League not detected as already finished"

    # The CLIENT computes the data needed to submit as an UPDATER: initStates, statesAtMatchday, scores.
    initPlayerStates = getInitPlayerStates(leagueIdx, ST_CLIENT)

    statesAtMatchday, scores = computeAllMatchdayStates(
        ST_CLIENT.leagues[leagueIdx].blockInit,
        ST_CLIENT.leagues[leagueIdx].blockStep,
        initPlayerStates,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        ST_CLIENT.leagues[leagueIdx].usersAlongData,
    )

    # ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the statesAtMatchday!:
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not-yet updated"
    initStatesHash          = serialHash(initPlayerStates)
    statesAtMatchdayHashes  = [serialHash(s) for s in statesAtMatchday]

    statesAtMatchdayHashesLie     = duplicate(statesAtMatchdayHashes)
    statesAtMatchdayHashesLie[0] += 1  # he lies about matchday 0 only

    ST.leagues[leagueIdx].updateLeague(
        initStatesHash,
        statesAtMatchdayHashesLie,
        scores,
        ADDR2,
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already challenged"

    # The CLIENT updates the league WITHOUT lying,
    # and additionally, stores the league pre-hash data, and updates every player involved
    ST_CLIENT.leagues[leagueIdx].updateLeague(
        initStatesHash,
        statesAtMatchdayHashes,
        scores,
        ADDR2,
        ST_CLIENT.currentBlock
    )
    updateClientAtEndOfLeague(leagueIdx, initPlayerStates, statesAtMatchday, scores, ST_CLIENT)
    assert ST_CLIENT.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already challenged"

    # A CHALLENGER tries to prove that the UPDATER lied with statesAtMatchday for matchday 0
    ST.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST_CLIENT.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    assert not ST.leagues[leagueIdx].isFullyVerified(ST.currentBlock)
    selectedMatchday    = 0
    prevMatchdayStates  = initPlayerStates  if selectedMatchday == 0 \
                                            else ST_CLIENT.leagues[leagueIdx].statesAtMatchday[selectedMatchday-1]

    ST.leagues[leagueIdx].challengeMatchdayStates(
        selectedMatchday,
        prevMatchdayStates,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx].usersAlongData),
        ST.currentBlock
    )
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful challenge"

    # ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the initStates!:
    ST.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST_CLIENT.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not updated"
    initStatesHashLie = duplicate(initStatesHash)+1

    ST.leagues[leagueIdx].updateLeague(
        initStatesHashLie,
        statesAtMatchdayHashes,
        scores,
        ADDR2,
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"


    # A CHALLENGER tries to prove that the UPDATER lied with the initHash
    ST.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST_CLIENT.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    assert not ST.leagues[leagueIdx].isFullyVerified(ST.currentBlock), "League not detected as not-yet fully verified"

    # create all the state of the environment (player team, owner, previous team ... league state)
    dataToChallengeInitStates = prepareDataToChallengeInitStates(leagueIdx, ST_CLIENT)

    ST.leagues[leagueIdx].challengeInitStates(
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(dataToChallengeInitStates),
        ST,
        ST.currentBlock,
        leagueIdx
    )
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful initHash challenge"


    # A nicer UPDATER now tells the truth:
    ST.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST_CLIENT.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST.leagues[leagueIdx].updateLeague(
        initStatesHash,
        statesAtMatchdayHashes,
        scores,
        ADDR2,
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    # ...and the CHALLENGER fails to prove anything
    ST.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST_CLIENT.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    selectedMatchday = 0
    ST.leagues[leagueIdx].challengeMatchdayStates(
        selectedMatchday,
        prevMatchdayStates,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx].usersAlongData),
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    ST.leagues[leagueIdx].challengeInitStates(
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(dataToChallengeInitStates),
        ST,
        ST.currentBlock,
        leagueIdx
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"


    # We do not wait enough and try to:
    #   create another league. It fails to do so because teams are still busy
    ST.advanceNBlocks(2)
    ST_CLIENT.advanceNBlocks(2)
    assert not ST.leagues[leagueIdx].isFullyVerified(ST.currentBlock), "League not detected as not-yet fully verified"
    blockInit = ST.currentBlock + 30
    usersInitData = {
        "teamIdxs": [teamIdx1, teamIdx3, teamIdx2, teamIdx4],
        "teamOrders": [DEFAULT_ORDER, DEFAULT_ORDER, REVERSE_ORDER, REVERSE_ORDER],
        "tactics": [TACTICS["433"], TACTICS["442"], TACTICS["433"], TACTICS["442"]]
    }
    try:
        leagueIdx2 = createLeague(ST.currentBlock, blockInit, blockStep, usersInitData, ST)
        itFailed = False
    except:
        itFailed = True

    assert itFailed, "League was created but previous league was not fully verified yet"


    # We do not wait enough and try to:
    #   sell/buy action is attempted, but fails because league is not full verified
    try:
        exchangePlayers(
            getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1, ST), ADDR1,
            getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6, ST), ADDR3,
            blockCounter.currentBlock,
            ST
        )
        itFailed = False
    except:
        itFailed = True

    assert itFailed, "Player sell/buy was allowed but previous league was not fully verified yet"


    # after waiting enough, the league gets fully verified and the new league can be created
    # ...with a player exchange just before the creation
    ST.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    ST_CLIENT.advanceNBlocks(CHALLENGING_PERIOD_BLKS-5)
    assert ST.leagues[leagueIdx].isFullyVerified(ST.currentBlock), "League not detected as already fully verified"

    playerIdx1 = getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1, ST)
    playerIdx2  = getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6, ST)

    exchangePlayers(
        playerIdx1, ADDR1,
        playerIdx2, ADDR3,
        ST
    )
    exchangePlayers(
        playerIdx1, ADDR1,
        playerIdx2, ADDR3,
        ST_CLIENT
    )
    assert getTeamIdxAndShirtForPlayerIdx(playerIdx1, ST) == (teamIdx4,6), "Exchange did not register properly in BC"
    assert getTeamIdxAndShirtForPlayerIdx(playerIdx2, ST) == (teamIdx1,1), "Exchange did not register properly in BC"

    # After the player exchange...
    leagueIdx2          = createLeague(blockInit, blockStep, usersInitData, ST)
    leagueIdx2_client   = createLeagueClient(blockInit, blockStep, usersInitData, ST_CLIENT)

    assert leagueIdx2 == leagueIdx2_client, "Leagues in client not in sync with BC"

    # An UPDATER updates:
    ST.advanceNBlocks(1000)
    ST_CLIENT.advanceNBlocks(1000)
    assert ST.leagues[leagueIdx].hasLeagueFinished(ST.currentBlock), "League should be finished by now"

    initPlayerStates = getInitPlayerStates(leagueIdx2, ST_CLIENT)

    statesAtMatchday, scores = computeAllMatchdayStates(
        ST.leagues[leagueIdx2].blockInit,
        ST.leagues[leagueIdx2].blockStep,
        initPlayerStates,
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersInitData),
        ST_CLIENT.leagues[leagueIdx2].usersAlongData,
    )

    # Having finished the league, the CLIENT can already update the CLIENT state
    updateClientAtEndOfLeague(leagueIdx2, initPlayerStates, statesAtMatchday, scores, ST_CLIENT)

    assert not ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "League not detected as not yet updated"
    initStatesHash = serialHash(initPlayerStates)
    statesAtMatchdayHashes = [serialHash(state) for state in statesAtMatchday]

    # compressedLeagueState encapsulates these hashes (in v2 maybe 1 single hash)
    #
    ST.leagues[leagueIdx2].updateLeague(
        initStatesHash,
        statesAtMatchdayHashes,
        scores,
        ADDR2,
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "League not detected as updated"

    # The CLIENT updates the league too,
    # and additionally, stores the league pre-hash data, and updates every player involved
    ST_CLIENT.leagues[leagueIdx2].updateLeague(
        initStatesHash,
        statesAtMatchdayHashes,
        scores,
        ADDR2,
        ST_CLIENT.currentBlock
    )
    updateClientAtEndOfLeague(leagueIdx2, initPlayerStates, statesAtMatchday, scores, ST_CLIENT)
    assert ST_CLIENT.leagues[leagueIdx2].hasLeagueBeenUpdated(), "League not detected as already challenged"

    # a challenger fails to prove that anything was wrong...
    # ... not for matchady = 0
    selectedMatchday    = 0
    prevMatchdayStates  = initPlayerStates  if selectedMatchday == 0 \
                                            else ST_CLIENT.leagues[leagueIdx2].statesAtMatchday[selectedMatchday-1]

    ST.leagues[leagueIdx2].challengeMatchdayStates(
        selectedMatchday,
        prevMatchdayStates,
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersAlongData),
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "Challenger was successful... but he should not be"

    # ... not for matchdat = 2
    selectedMatchday    = 2
    prevMatchdayStates  = initPlayerStates  if selectedMatchday == 0 \
                                            else ST_CLIENT.leagues[leagueIdx2].statesAtMatchday[selectedMatchday-1]

    ST.leagues[leagueIdx2].challengeMatchdayStates(
        selectedMatchday,
        prevMatchdayStates,
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersAlongData),
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "Challenger was successful... but he should not be"

    # ... not for initStates
    dataToChallengeInitStates = prepareDataToChallengeInitStates(leagueIdx2, ST_CLIENT)

    ST.leagues[leagueIdx2].challengeInitStates(
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersInitData),
        duplicate(dataToChallengeInitStates),
        ST,
        ST.currentBlock,
        leagueIdx2
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "Challenger was successful... but he should not be"

    # We make sure that we can inquire the state of any player after these leagues and player sales:
    player1State = getLastWrittenPlayerStateFromPlayerIdx(1, ST_CLIENT)
    player1ChallengeData = computeDataToChallengePlayerIdx(1, ST_CLIENT)
    assert isCorrectStateForPlayerIdx(player1State, player1ChallengeData, ST), "Computed player state by CLIENT is not recognized by BC.."

    # The following all-team printout is interesting. On the one hand, it checks that all player states
    # in that team can be certified by the BC. On the other hand, you can check that the 2nd player
    # corresponds to the player bought from team4, in the exchange done above.
    printTeam(teamIdx1, ST_CLIENT)

    # Returns test result, to later check against expected
    testResult = intHash(serialize(ST) + serialize(ST_CLIENT)) % 1000
    return testResult


def runTest(name, result, expected):
    success = (result == expected)
    if success:
        print name + ": PASSED"
    else:
        print name + ": FAILED.  Result(%s) vs Expected(%s) " % (result, expected)
    return success


success = True
success = success and runTest(name = "Test Simple Team Creation", result = test1(), expected = 10222)
success = success and runTest(name = "Test Entire Workflow",      result = test2(), expected = 145)
if success:
    print "ALL TESTS:  -- PASSED --"
else:
    print "At least one test FAILED"