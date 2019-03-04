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

    print("Team created with teamIdx, teamName = " + str(teamIdx1) + ", " + ST.teams[teamIdx1].name)
    hash0 = printTeam(teamIdx1, ST_CLIENT)

    print("\n\nplayers 2 and 24 before sale:\n")
    hash1 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(2, ST_CLIENT))

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "PlayerStates not in sync BC vs client"

    print("\n")
    hash2 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(24, ST_CLIENT))

    advanceNBlocks(10, ST, ST_CLIENT)

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

    print("\n\nplayers 2 and 24 after sale:\n")
    hash3 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(2, ST))
    print("\n")
    hash4 = printPlayer(getLastWrittenPlayerStateFromPlayerIdx(24, ST_CLIENT))
    hashSum         = hash0+hash1+hash2+hash3+hash4
    return hashSum



# TEST2: all the workflow!
def test2():
    # Create contract storage in BC, and its extended version for the CLIENT
    # We need to keep both in sync. The CLIENT stores, besides what is in the BC, the pre-hash stuff.
    ST          = Storage()
    ST_CLIENT   = Storage()
    ST_CLIENT.addAccumulator()

    advanceToBlock(10, ST, ST_CLIENT)

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
    advanceToBlock(100, ST, ST_CLIENT)

    # One verse is about 1 hour, so a day is about 24 verseSteps
    verseInit = 3
    verseStep = 24

    usersInitData = {
        "teamIdxs": [teamIdx1, teamIdx2],
        "teamOrders": [DEFAULT_ORDER, REVERSE_ORDER],
        "tactics": [TACTICS["442"], TACTICS["433"]]
    }


    # Create league in BC and CLIENT. The latter stores things pre-hash too
    leagueIdx          = createLeague(verseInit, verseStep, usersInitData, ST)
    leagueIdx_client   = createLeagueClient(verseInit, verseStep, usersInitData, ST_CLIENT)

    assert ST.leagues[leagueIdx].isLeagueIsAboutToStart(ST.currentVerse), "League not detected as created"

    # advance a bit before first match to change tactics
    assert ST.currentVerse == 0, "We should start with verse 0"
    advanceToBlock(ST.nextVerseBlock()-5, ST, ST_CLIENT)

    action00 = {"teamIdx": teamIdx1, "teamOrder": ORDER1, "tactics": TACTICS["433"]}
    ST_CLIENT.accumulateAction(action00)


    # Advance to just before matchday 2, which starts at verse 3 + 24 = 27
    # From verse 0 to 26:
    assert ST.currentVerse == 0, "We should start with verse 0"
    advanceNVerses(24, ST, ST_CLIENT)
    assert ST.currentVerse == 24, "We should be at verse 24, league finishes at 27"
    advanceToBlock(ST.nextVerseBlock()-5, ST, ST_CLIENT)

    assert ST.leagues[leagueIdx].hasLeagueStarted(ST.currentVerse), "League not detected as already being played"
    assert not ST.leagues[leagueIdx].hasLeagueFinished(ST.currentVerse), "League not detected as not finished yet"

    # Cook data to change tactics before games in matchday 2 begin.
    action0 = {"teamIdx": teamIdx1, "teamOrder": ORDER2, "tactics": TACTICS["433"]}
    action1 = {"teamIdx": teamIdx2, "teamOrder": DEFAULT_ORDER, "tactics": TACTICS["442"]}

    ST_CLIENT.accumulateAction(action0)

    advanceNVerses(2, ST, ST_CLIENT)
    assert ST.currentVerse == 26, "We should be at verse 26, league finishes at 27"
    advanceToBlock(ST.nextVerseBlock()-5, ST, ST_CLIENT)
    ST_CLIENT.accumulateAction(action1)



    # Move beyond league end
    advanceNVerses(1, ST, ST_CLIENT)
    assert ST.leagues[leagueIdx].hasLeagueFinished(ST.currentVerse), "League not detected as already finished"

    initPlayerStates = ST_CLIENT.getInitPlayerStates(leagueIdx)
    statesAtMatchday, tacticsAtMatchDay, teamOrdersAtMatchDay, scores = ST_CLIENT.computeAllMatchdayStates(leagueIdx)


    # TODO: treat initStates the same way as states and avoid initPlayerHash being different

    # ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the statesAtMatchday!:
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not-yet updated"
    initStatesHash          = serialHash(initPlayerStates)
    dataAtMatchdayHashes = computesdataAtMatchdayHashes(statesAtMatchday, tacticsAtMatchDay, teamOrdersAtMatchDay)

    dataAtMatchdayHashesLie     = duplicate(dataAtMatchdayHashes)
    dataAtMatchdayHashesLie[0] += 1  # he lies about matchday 0 only

    ST.updateLeague(
        leagueIdx,
        initStatesHash,
        dataAtMatchdayHashesLie,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"

    # The CLIENT updates the league WITHOUT lying,
    # and additionally, stores the league pre-hash data, and updates every player involved
    ST_CLIENT.updateLeague(
        leagueIdx,
        initStatesHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    updateClientAtEndOfLeague(leagueIdx, initPlayerStates, statesAtMatchday, tacticsAtMatchDay, teamOrdersAtMatchDay, scores, ST_CLIENT)
    assert ST_CLIENT.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already challenged"

    # A CHALLENGER tries to prove that the UPDATER lied with statesAtMatchday for matchday 0
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.leagues[leagueIdx].isFullyVerified(ST.currentBlock)
    selectedMatchday    = 0
    prevMatchdayStates, prevMatchdayTactics, prevMatchdayTeamOrders = \
        getPrevMatchdayData(ST_CLIENT, leagueIdx, selectedMatchday)

    merkleProof, values, depth = ST_CLIENT.getMerkleProof(leagueIdx, selectedMatchday)


    ST.challengeMatchdayStates(
        leagueIdx,
        selectedMatchday,
        prevMatchdayStates,
        prevMatchdayTactics,
        prevMatchdayTeamOrders,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(allActionsInThisLeague[selectedMatchday]),
        merkleProof,
        values,
        depth
    )


    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful challenge"

    # ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the initStates!:
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
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
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.leagues[leagueIdx].isFullyVerified(ST.currentBlock), "League not detected as not-yet fully verified"

    dataToChallengeInitStates = prepareDataToChallengeInitStates(leagueIdx, ST_CLIENT)

    ST.leagues[leagueIdx].challengeInitStates(
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(dataToChallengeInitStates),
        ST,
        ST.currentBlock
    )
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful initHash challenge"


    # A nicer UPDATER now tells the truth:
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    ST.leagues[leagueIdx].updateLeague(
        initStatesHash,
        statesAtMatchdayHashes,
        scores,
        ADDR2,
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    # ...and the CHALLENGER fails to prove anything
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    selectedTeam = 0
    ST.leagues[leagueIdx].challengeMatchdayStates(
        selectedTeam,
        initPlayerStates,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx].usersAlongData),
        ST.currentBlock
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    # We do not wait enough and try to:
    #   create another league. It fails to do so because teams are still busy
    advanceNBlocks(2, ST, ST_CLIENT)
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
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
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
    advanceNBlocks(1000, ST, ST_CLIENT)
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

    # We make sure that we can inquire the state of any player after these leagues and player sales:
    player1State = getLastWrittenPlayerStateFromPlayerIdx(1, ST_CLIENT)
    player1ChallengeData = computeDataToChallengePlayerIdx(1, ST_CLIENT)
    assert isCorrectStateForPlayerIdx(player1State, player1ChallengeData, ST), "Computed player state by CLIENT is not recognized by BC.."

    # The following all-team printout is interesting. On the one hand, it checks that all player states
    # in that team can be certified by the BC. On the other hand, you can check that the 2nd player
    # corresponds to the player bought from team4, in the exchange done above.
    printTeam(teamIdx1, ST_CLIENT)

    # Returns test result, to later check against expected
    testResult = intHash(serialize2str(ST) + serialize2str(ST_CLIENT)) % 1000
    return testResult

def test3():
    ST          = Storage()
    ST_CLIENT   = Storage()
    ST_CLIENT.addAccumulator()

    advanceToBlock(10, ST, ST_CLIENT)

    action0 = {"teamIdx": 3, "tactics": TACTICS["433"]}
    action1 = {"teamIdx": 2, "tactics": TACTICS["442"]}
    action3 = {"teamIdx": 22, "tactics": TACTICS["433pressing"]}
    action4 = {"teamIdx": 33, "tactics": TACTICS["442pressing"]}

    ST_CLIENT.accumulateAction(action0)
    ST_CLIENT.accumulateAction(action1)

    assert ST.currentVerse == ST_CLIENT.currentVerse == 0, "Error: Starting verse should be 0 (dummy verse)."
    assert len(ST_CLIENT.VerseCommits)==1, "Error: CLIENT should start with a single dummy commit."
    assert len(ST.VerseCommits)==1, "Error: BC should start with a single dummy commit."
    assert len(ST_CLIENT.Accumulator.buffer)==1, "Error: CLIENT should have accumulated actions from 1 block"
    assert len(ST_CLIENT.Accumulator.buffer[10])==2, "Error: CLIENT should have accumulated 2 actions for 1st block"

    advanceNBlocks(3, ST, ST_CLIENT)

    ST_CLIENT.accumulateAction(action3)
    assert len(ST_CLIENT.Accumulator.buffer[13])==1, "Error: CLIENT should have accumulated another action"

    advanceNBlocks(360, ST, ST_CLIENT)

    assert ST.currentVerse == ST_CLIENT.currentVerse == 1, "Error: Current verse is outdated"
    assert len(ST_CLIENT.VerseCommits)==2, "Error: CLIENT should have made the first non-dummy commit already."
    assert len(ST.VerseCommits)==2, "Error: BC should have made the first non-dummy commit already."
    assert len(ST_CLIENT.Accumulator.buffer)==0, "Error: CLIENT should have emptied the actions buffer"
    assert len(ST_CLIENT.Accumulator.commitedActions)==2, "Error: CLIENT should have updated the commited actions"


    ST_CLIENT.accumulateAction(action4)


    testResult = intHash(serialize2str(ST) + serialize2str(ST_CLIENT)) % 1000
    return testResult


def test4():
    leafs = [1,2,3,4,5,6,"rew"]
    tree, depth = make_tree(leafs, serialHash)
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
    tree, depth = make_tree(leafs, serialHash)
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
# success = success and runTest(name = "Test Simple Team Creation", result = test1(), expected = 9207)
success = success and runTest(name = "Test Entire Workflow",      result = test2(), expected = 935)
# success = success and runTest(name = "Test Accumulator",      result = test3(), expected = 396)
# success = success and runTest(name = "Test Merkle",      result = test4(), expected = True)
if success:
    print("ALL TESTS:  -- PASSED --")
else:
    print("At least one test FAILED")



# TODO:
# - change teamIdx for teamPos inside usersAlongData
# - gather dataAtMAtch day as a struct
#   - likeweise, put initStates as states at 0 (not sure)
# gather all merkle proof data (vals, hashes, depth) in a struct