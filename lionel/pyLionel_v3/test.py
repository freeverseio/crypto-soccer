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
    ST          = Storage(isClient = False)
    ST_CLIENT   = Storage(isClient = True)

    teamIdx1 = ST.createTeam("Barca", ADDR1)
    teamIdx2 = ST.createTeam("Madrid", ADDR2)

    teamIdx1_client = ST_CLIENT.createTeam("Barca", ADDR1)
    teamIdx2_client = ST_CLIENT.createTeam("Madrid", ADDR2)

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "TeamIdx not in sync BC vs client"

    # Test that we can ask the BC if state of a player (computed by the Client) is correct:
    player1State                = ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(1)
    dataToChallengePlayerState  = ST_CLIENT.computeDataToChallengePlayerIdx(1)
    assert ST.isCorrectStateForPlayerIdx(player1State, dataToChallengePlayerState), "Computed player state by CLIENT is not recognized by BC.."

    print("Team created with teamIdx, teamName = " + str(teamIdx1) + ", " + ST.teams[teamIdx1].name)
    hash0 = printTeam(teamIdx1, ST_CLIENT)

    print("\n\nplayers 2 and 24 before sale:\n")
    hash1 = printPlayer(ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(2))

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "PlayerStates not in sync BC vs client"

    print("\n")
    hash2 = printPlayer(ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(24))

    advanceNBlocks(10, ST, ST_CLIENT)

    ST.exchangePlayers(
        2, ADDR1,
        24, ADDR2
    )
    ST_CLIENT.exchangePlayers(
        2, ADDR1,
        24, ADDR2
    )

    print("\n\nplayers 2 and 24 after sale:\n")
    hash3 = printPlayer(ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(2))
    print("\n")
    hash4 = printPlayer(ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(24))
    hashSum         = hash0+hash1+hash2+hash3+hash4
    return hashSum



# TEST2: all the workflow!
def test2():
    # Create contract storage in BC, and its extended version for the CLIENT
    # We need to keep both in sync. The CLIENT stores, besides what is in the BC, the pre-hash stuff.
    ST          = Storage(isClient = False)
    ST_CLIENT   = Storage(isClient = True)

    # The accumulator is responsible for receving user actions and committing them in the correct verse.
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
    ST_CLIENT.accumulateAction(action00)


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

    # Move beyond league end
    advanceNVerses(1, ST, ST_CLIENT)
    assert ST.hasLeagueFinished(leagueIdx), "League not detected as already finished"

    # CLIENT computes the data needed to update league
    initPlayerStates        = ST_CLIENT.getInitPlayerStates(leagueIdx)
    dataAtMatchdays, scores = ST_CLIENT.computeAllMatchdayStates(leagueIdx)

    # ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the statesAtMatchday!:
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not-yet updated"

    initStatesHash       = serialHash(initPlayerStates)

    dataAtMatchdayHashes, lastDayTree = ST_CLIENT.prepareHashesForDataAtMatchdays(dataAtMatchdays)

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

    # The CLIENT updates the league WITHOUT lying...
    ST_CLIENT.updateLeague(
        leagueIdx,
        initStatesHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    # ...and additionally, stores the league pre-hash data, and updates every player involved
    ST_CLIENT.storePreHashDataInClientAtEndOfLeague(leagueIdx, initPlayerStates, dataAtMatchdays, lastDayTree, scores)
    assert ST_CLIENT.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"

    # A CHALLENGER tries to prove that the UPDATER lied with statesAtMatchday for matchday 0
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.isFullyVerified(leagueIdx)

    # ...first, it selects a matchday, and gathers the data at that matchday (states, tactics, teamOrders)
    selectedMatchday    = 0
    dataAtPrevMatchday = ST_CLIENT.getPrevMatchdayData(leagueIdx, selectedMatchday)

    # ...next, it builds the Merkle proof for the actions commited on the corresponding verse, for that league
    merkleProofDataForMatchday = ST_CLIENT.getMerkleProof(leagueIdx, selectedMatchday)

    # ...finally, it does the challenge
    ST.challengeMatchdayStates(
        leagueIdx,
        selectedMatchday,
        dataAtPrevMatchday,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx].actionsPerMatchday[selectedMatchday]),
        merkleProofDataForMatchday
    )
    # Since it must succeed, the league is 'reset', without any update
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful challenge"


    # More lies: the CLIENT, acting as an UPDATER, submits to the BC... a lie in the initStates!:
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not updated"
    initStatesHashLie = duplicate(initStatesHash)+1

    ST.updateLeague(
        leagueIdx,
        initStatesHashLie,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"

    # A CHALLENGER proves that the UPDATER lied with the initHash
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.isFullyVerified(leagueIdx), "League not detected as not-yet fully verified"

    # ...first it gathers the data needed to challenge the init states
    dataToChallengeInitStates = ST_CLIENT.prepareDataToChallengeInitStates(leagueIdx)
    ST.challengeInitStates(
        leagueIdx,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(dataToChallengeInitStates),
    )
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful initHash challenge"


    # Finally, some truth: a nicer UPDATER now tells the truth:
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    ST.updateLeague(
        leagueIdx,
        initStatesHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    # ...and the CHALLENGER fails to prove anything
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    ST.challengeInitStates(
        leagueIdx,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(dataToChallengeInitStates),
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    # We do not wait enough and try to create another league with these teams.
    # It fails to do so because teams are still busy
    advanceNBlocks(2, ST, ST_CLIENT)
    assert not ST.isFullyVerified(leagueIdx), "League not detected as not-yet fully verified"
    verseInit = ST.currentVerse + 30
    usersInitData = {
        "teamIdxs": [teamIdx1, teamIdx3, teamIdx2, teamIdx4],
        "teamOrders": [DEFAULT_ORDER, DEFAULT_ORDER, REVERSE_ORDER, REVERSE_ORDER],
        "tactics": [TACTICS["433"], TACTICS["442"], TACTICS["433"], TACTICS["442"]]
    }
    try:
        leagueIdx2 = ST.createLeague(verseInit, verseStep, usersInitData)
        itFailed = False
    except:
        itFailed = True

    assert itFailed, "League was created but previous league was not fully verified yet"


    # We do not wait enough and try to:
    #   sell/buy action is attempted, but fails because league is not full verified
    try:
        ST.exchangePlayers(
            ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx1, 1), ADDR1,
            ST.getPlayerIdxFromTeamIdxAndShirt(teamIdx4, 6), ADDR3
        )
        itFailed = False
    except:
        itFailed = True

    assert itFailed, "Player sell/buy was allowed but previous league was not fully verified yet"


    # after waiting enough, the league gets fully verified and the new league can be created
    # ...with a player exchange just before the creation
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert ST.isFullyVerified(leagueIdx), "League not detected as already fully verified"

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
    assert ST.getTeamIdxAndShirtForPlayerIdx(playerIdx1) == (teamIdx4,6), "Exchange did not register properly in BC"
    assert ST.getTeamIdxAndShirtForPlayerIdx(playerIdx2) == (teamIdx1,1), "Exchange did not register properly in BC"

    # After the player exchange... we create another league
    leagueIdx2          = ST.createLeague(verseInit, verseStep, usersInitData)
    leagueIdx2_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)

    assert leagueIdx2 == leagueIdx2_client, "Leagues in client not in sync with BC"

    # At the end of league, an UPDATER updates telling the truth:
    advanceNVerses(1000, ST, ST_CLIENT)
    assert ST.hasLeagueFinished(leagueIdx2), "League should be finished by now"

    initPlayerStates = ST_CLIENT.getInitPlayerStates(leagueIdx2)
    dataAtMatchdays, scores = ST_CLIENT.computeAllMatchdayStates(leagueIdx2)

    initStatesHash       = serialHash(initPlayerStates)
    dataAtMatchdayHashes, lastDayTree = ST_CLIENT.prepareHashesForDataAtMatchdays(dataAtMatchdays)
    assert dataAtMatchdayHashes[0] == serialHash(dataAtMatchdays[0]), "Something went wrong preparing hashes"
    ST.updateLeague(
        leagueIdx2,
        initStatesHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "League not detected as already updated"

    # The CLIENT updates the league too,
    # and additionally, stores the league pre-hash data, and updates every player involved
    ST_CLIENT.updateLeague(
        leagueIdx2,
        initStatesHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    ST_CLIENT.storePreHashDataInClientAtEndOfLeague(leagueIdx2, initPlayerStates, dataAtMatchdays, lastDayTree, scores)
    assert ST_CLIENT.leagues[leagueIdx2].hasLeagueBeenUpdated(), "League not detected as already updated"

    # A challenger fails to prove anything is wrong with init states...
    dataToChallengeInitStates = ST_CLIENT.prepareDataToChallengeInitStates(leagueIdx2)
    ST.challengeInitStates(
        leagueIdx2,
        ST_CLIENT.leagues[leagueIdx2].usersInitData,
        duplicate(dataToChallengeInitStates),
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"

    # ...or with matchday 0...
    selectedMatchday = 0
    dataAtPrevMatchday = ST_CLIENT.getPrevMatchdayData(leagueIdx2, selectedMatchday)
    merkleProofDataForMatchday = ST_CLIENT.getMerkleProof(leagueIdx2, selectedMatchday)
    ST.challengeMatchdayStates(
        leagueIdx2,
        selectedMatchday,
        dataAtPrevMatchday,
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx2].actionsPerMatchday[selectedMatchday]),
        merkleProofDataForMatchday
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"

    # ...or with matchday 4...
    selectedMatchday = 5
    dataAtPrevMatchday = ST_CLIENT.getPrevMatchdayData(leagueIdx2, selectedMatchday)
    merkleProofDataForMatchday = ST_CLIENT.getMerkleProof(leagueIdx2, selectedMatchday)
    ST.challengeMatchdayStates(
        leagueIdx2,
        selectedMatchday,
        dataAtPrevMatchday,
        duplicate(ST_CLIENT.leagues[leagueIdx2].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx2].actionsPerMatchday[selectedMatchday]),
        merkleProofDataForMatchday
    )
    assert ST.leagues[leagueIdx2].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"


    # We make sure that we can inquire the state of any player after these leagues and player sales:
    player1State = ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(1)
    dataToChallengePlayerState = ST_CLIENT.computeDataToChallengePlayerIdx(1)
    assert ST.isCorrectStateForPlayerIdx(player1State, dataToChallengePlayerState), "Computed player state by CLIENT is not recognized by BC.."

    # The following all-team printout is interesting. On the one hand, it checks that all player states
    # in that team can be certified by the BC. On the other hand, you can check that the 2nd player
    # corresponds to the player bought from team4, in the exchange done above.
    printTeam(teamIdx1, ST_CLIENT)

    # create many teams, and leagues, and mess it all.
    advanceNVerses(1000, ST, ST_CLIENT)
    teamIdxs = []
    for t in range(100):
        teamIdxs.append(ST.createTeam("BotTeam"+str(t), ADDR1))
        ST_CLIENT.createTeam("BotTeam"+str(t), ADDR2)

    for p in range(400):
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
        print(p)
        playerState = ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(playerIdx1)
        dataToChallengePlayerState = ST_CLIENT.computeDataToChallengePlayerIdx(playerIdx1)
        assert ST.isCorrectStateForPlayerIdx(playerState, dataToChallengePlayerState), "Computed player state by CLIENT is not recognized by BC.."



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
success = success and runTest(name = "Test Simple Team Creation", result = test1(), expected = 9207)
success = success and runTest(name = "Test Entire Workflow",      result = test2(), expected = 196)
# success = success and runTest(name = "Test Accumulator",      result = test3(), expected = 396)
success = success and runTest(name = "Test Merkle",      result = test4(), expected = True)
if success:
    print("ALL TESTS:  -- PASSED --")
else:
    print("At least one test FAILED")



# TODO:
# BUG: getLastWrittenPlayerStateFromPlayerIdx does not really return last written state in BC, but
#  last written in Client
#   - likeweise, put initStates as states at 0 (not sure)
# treat initStates the same way as states and avoid initPlayerHash being different
#         # TODO: check that the provided state proofs contain the actual player idx!!!!! --> see structs challengeinit hash
# add test for multiple simultaneous leauges (for the proof), some with actions, some without, etc
# use merkle proof for playerStates at previous league?
# how can we store the hash of the teamIdx???? can we not sign a team in another league?
#     I think that the answer is that we store the league in the team property!
# gather together code to update actions, e.g., find all  not actionsAtSelectedMatchday == 0

# Note that dataAtMatchday.states = after the given match!
# remove ugly:         if type(dataToChallengePlayerState) == type(DataAtMatchday(0, 0, 0)):

# leafIdx = list(dataToChallengePlayerState.values.keys())[0]
# isPlayerStateInsideDataToChallenge => not need anymore, right? it's inside getPlayerStateFromChallengeData already
# test getOwner works and use it in player exchange tests

# TODO: - less important -
# do not store scores but the hash or merkle root
# unify all iniHash, serialHash, etc
# remove need for the last matchdayHash, as we just need to test the states.
