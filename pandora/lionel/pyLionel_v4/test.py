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
    player1State                = ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(1)
    dataToChallengePlayerState  = ST_CLIENT.computeDataToChallengePlayerSkills(1)
    assert ST.areLatestSkills(player1State, dataToChallengePlayerState), "Computed player state by CLIENT is not recognized by BC.."

    print("Team created with teamIdx, teamName = " + str(teamIdx1) + ", " + ST.teams[teamIdx1].name)
    hash0 = printTeam(teamIdx1, ST_CLIENT)

    print("\n\nplayers 2 and 24 before sale:\n")

    hash1 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(2))

    assert (teamIdx1 == teamIdx1_client) and (teamIdx2 == teamIdx2_client), "PlayerStates not in sync BC vs client"

    print("\n")
    hash2 = printPlayerFromSkills(ST_CLIENT, ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(24))

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
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not-yet updated"

    # CLIENT computes the data needed to update league (and stores it in the CLIENT)
    initSkillsHash, dataAtMatchdayHashes, scores = ST_CLIENT.updateLeagueInClient(leagueIdx, ADDR2)

    # ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the skillsAtMatchday!:
    dataAtMatchdayHashesLie     = duplicate(dataAtMatchdayHashes)
    dataAtMatchdayHashesLie[0] += 1  # he lies about matchday 0 only

    ST.updateLeague(
        leagueIdx,
        initSkillsHash,
        dataAtMatchdayHashesLie,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"

    # A CHALLENGER tries to prove that the UPDATER lied with skillsAtMatchday for matchday 0
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)

    selectedMatchday = 0
    challengeLeagueAtSelectedMatchday(selectedMatchday, leagueIdx, ST, ST_CLIENT)
    # Since it must succeed, the league is 'reset', without any update
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful challenge"


    # More lies: the CLIENT, acting as an UPDATER, submits to the BC... a lie in the initStates!:
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not updated"
    initSkillsHashLie = duplicate(initSkillsHash)+1

    ST.updateLeague(
        leagueIdx,
        initSkillsHashLie,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"

    # A CHALLENGER proves that the UPDATER lied with the initHash
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    assert not ST.isFullyVerified(leagueIdx), "League not detected as not-yet fully verified"

    # ...first it gathers the data needed to challenge the init states
    ST.challengeInitSkills(
        leagueIdx,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
    )
    assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not reset after successful initHash challenge"


    # Finally, some truth: a nicer UPDATER now tells the truth:
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    ST.updateLeague(
        leagueIdx,
        initSkillsHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as updated"

    # ...and the CHALLENGER fails to prove anything
    advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5, ST, ST_CLIENT)
    ST.challengeInitSkills(
        leagueIdx,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
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
        leagueIdx = ST.createLeague(verseInit, verseStep, usersInitData)
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
    leagueIdx          = ST.createLeague(verseInit, verseStep, usersInitData)
    leagueIdx_client   = ST_CLIENT.createLeagueClient(verseInit, verseStep, usersInitData)

    assert leagueIdx == leagueIdx_client, "Leagues in client not in sync with BC"

    # At the end of league, an UPDATER updates telling the truth:
    advanceNVerses(1000, ST, ST_CLIENT)
    assert ST.hasLeagueFinished(leagueIdx), "League should be finished by now"

    initSkillsHash, dataAtMatchdayHashes, scores = ST_CLIENT.updateLeagueInClient(leagueIdx, ADDR2)
    print(initSkillsHash)
    print(pylio.serialHash(ST_CLIENT.leagues[leagueIdx].getInitPlayerSkills()))

    ST.updateLeague(
        leagueIdx,
        initSkillsHash,
        dataAtMatchdayHashes,
        scores,
        ADDR2,
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"


    # A challenger fails to prove anything is wrong with init states...
    ST.challengeInitSkills(
        leagueIdx,
        ST_CLIENT.leagues[leagueIdx].usersInitData,
        duplicate( ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills )
    )
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"

    # ...or with matchday 0...
    # challegeLeagueAtSelectedMatchday(leagueIdx)
    selectedMatchday = 0
    challengeLeagueAtSelectedMatchday(selectedMatchday, leagueIdx, ST, ST_CLIENT)
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"

    # ...or with matchday 4...
    selectedMatchday = 5
    challengeLeagueAtSelectedMatchday(selectedMatchday, leagueIdx, ST, ST_CLIENT)
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"


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
        ST_CLIENT.createTeam("BotTeam"+str(t), ADDR2)

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

        if l==0:
            firstLeagueIdx = duplicate(leagueIdx)

        assert (leagueIdx == leagueIdx_client), "leagueIdx not in sync BC vs client"
        assert ST.isLeagueIsAboutToStart(leagueIdx), "League not detected as created"
        advanceNVerses(intHash(str(l))%2 , ST, ST_CLIENT) # advance either 1 or 0 verses

    advanceNVerses(1000, ST, ST_CLIENT)
    nActionsPerLoop = 3
    for l in range(nLeagues):
        advanceNVerses(intHash(str(l))%27 , ST, ST_CLIENT) # advance any number of verses between 0, 27
        leagueIdx = firstLeagueIdx + l
        assert ST.hasLeagueFinished(leagueIdx), "League not detected as already finished"
        assert not ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as not-yet updated"

        advanceNVerses(1 , ST, ST_CLIENT)

        for a in range(nActionsPerLoop):
            thisTeamIdx = getRandomElement(ST_CLIENT.leagues[leagueIdx].usersInitData["teamIdxs"],l+a)
            action =  {
                "teamIdx": thisTeamIdx,
                "teamOrder": getRandomElement(POSSIBLE_ORDERS,l+a),
                "tactics": getRandomElement(POSSIBLE_TACTICS,l+a+13)
            }
            advanceNVerses(intHash(str(l+a+14))%2, ST, ST_CLIENT) # advance either 0 or 1 verse.
            ST_CLIENT.accumulateAction(action)

        initSkillsHash, dataAtMatchdayHashes, scores = ST_CLIENT.updateLeagueInClient(leagueIdx, ADDR2)

        ST.updateLeague(
            leagueIdx,
            initSkillsHash,
            dataAtMatchdayHashes,
            scores,
            ADDR2,
        )
        assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"

        # A challenger fails to prove anything is wrong with init states...
        ST.challengeInitSkills(
            leagueIdx,
            ST_CLIENT.leagues[leagueIdx].usersInitData,
            duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
        )
        assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"

        # ...or for any of the total number of matchdays
        nDays = len( ST.leagues[leagueIdx].dataAtMatchdayHashes)-1 # the last one is the merkle root
        for selectedMatchday in range(nDays):
            challengeLeagueAtSelectedMatchday(selectedMatchday, leagueIdx, ST, ST_CLIENT)
            assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "Challenger was successful when he should not be"


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
# success = success and runTest(name = "Test Simple Team Creation", result = test1(), expected = 10754)
success = success and runTest(name = "Test Entire Workflow",      result = test2(), expected = 353)
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