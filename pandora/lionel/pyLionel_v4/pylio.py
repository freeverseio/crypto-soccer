import numpy as np
import copy
import datetime
import sha3
from copy import deepcopy as duplicate
from constants import *
from structs import *
from merkle_tree import *


# ------------ Functions to take hashes, serialize structs ------------

#  serializes and converts to str in a complicated way
def serialize2str(object):
    return str(serialize(object).hex())

# Numpy accepts a max possible seed
def limitSeed(seed):
    return seed % MAX_RND_SEED_ALLOWED_BY_NUMPY

# Returns keccak of string in hex format
def hexHash(str):
    return sha3.keccak_256(str.encode('utf-8')).hexdigest()

# Returns kekkack of string in decimal format
def intHash(str):
    return int(hexHash(str), 16)

# serializes object and then taken intHash
def serialHash(obj):
    return intHash(serialize2str(obj))


# ------------ Functions to print structs ------------

# Simple player print
def printPlayer(playerState):
    toPrint =  "PlayerIdx: %s\n" % str(playerState.getPlayerIdx())
    toPrint += "Age      : %s\n" % str(playerState.getMonth()/12)
    toPrint += "Skills   : %s\n" % str(playerState.getSkills())
    toPrint += "TeamIdx  : %s\n" % str(playerState.getCurrentTeamIdx())
    toPrint += "ShirtNum : %s\n" % str(playerState.getCurrentShirtNum())
    toPrint += "SaleBlock: %s\n" % str(playerState.getLastSaleBlocknum())
    print("%s" % toPrint)
    return intHash(toPrint) % 1000

# Simple team print
def printTeam(teamIdx, ST_CLIENT):
    hash = 0
    print("Player for teamIdx %d, with teamName %s: " %(teamIdx, ST_CLIENT.teams[teamIdx].name))
    for shirtNum in range(NPLAYERS_PER_TEAM):
        playerIdx = ST_CLIENT.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
        playerSkills = ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(playerIdx)
        playerChallengeData = ST_CLIENT.computeDataToChallengePlayerSkills(playerSkills.getPlayerIdx())
        assert ST_CLIENT.areLatestSkills(playerSkills, playerChallengeData), "Player state not correctly in sync"
        hash += printPlayer(ST_CLIENT.skillsToLastWrittenState(playerSkills))
    return hash

def isValidOrdering(playerOrders):
    # TODO: Currently not implemented. Check all nums are different and in [0, NPLAYERS_PER_TEAM]
    return True

def shiftBack(t, nTeams):
    if (t < nTeams):
        return t
    else:
        return t-(nTeams-1)


def getTeamsInMatchFirstHalf(matchday, match, nTeams):
    team1 = 0
    if (match > 0):
        team1 = shiftBack(nTeams-match+matchday, nTeams)

    team2 = shiftBack(match+1+matchday, nTeams)
    if ( (matchday % 2) == 0):
        return team1, team2
    else:
        return team2, team1


def getTeamsInMatch(matchday, match, nTeams):
    assert matchday < 2 * (nTeams - 1), "This league does not have so many matchdays"
    if (matchday < (nTeams - 1)):
        (team1, team2) = getTeamsInMatchFirstHalf(matchday, match, nTeams)
    else:
        (team2, team1) = getTeamsInMatchFirstHalf(matchday - (nTeams - 1), match, nTeams);
    return team1, team2

# mockup
def playMatch(initPlayerSkills1, initPlayerSkills2, tactics1, tactics2, teamOrders1, teamOrders2, MatchSeed):
    hash1 = intHash(str(MatchSeed)+serialize2str(initPlayerSkills1)+serialize2str(tactics1)+serialize2str(teamOrders1))
    hash2 = intHash(str(MatchSeed)+serialize2str(initPlayerSkills2)+serialize2str(tactics2)+serialize2str(teamOrders2))
    return hash1 % 4, hash2 % 4

# mockup: the rating of a team is just the sum of all skills
def computeTeamRating(playerSkills):
    return sum([sum(thisPlayerSkills.getSkills()) for thisPlayerSkills in playerSkills])


def addFixedPointsToAllPlayers(playerSkills, points):
    for playerState in playerSkills:
        playerState.setSkills(playerState.getSkills() + points)


# given the result, it computes the evolution points won per team, and applies them to their players
def updatePlayerSkillsAfterMatch(playerSkills1, playerSkills2, goals1, goals2):
    ps1 = duplicate(playerSkills1)
    ps2 = duplicate(playerSkills2)

    if goals1 == goals2:
        return ps1, ps2

    pointsWon = computePointsWon(ps1, ps2, goals1, goals2)
    if goals1 > goals2:
        addFixedPointsToAllPlayers(ps1, pointsWon)
    else:
        addFixedPointsToAllPlayers(ps2, pointsWon)

    return ps1, ps2


# simple mockup of what the evolution points could look like.
def computePointsWon(playerState1, playerState2, goals1, goals2):
    ratingDiff              = computeTeamRating(playerState1) - computeTeamRating(playerState2)
    winnerWasBetter         = (ratingDiff > 0 and goals1>goals2) or (ratingDiff < 0 and goals1<goals2)

    if ratingDiff == 0:
        return 5
    else:
        return (2 if winnerWasBetter else 10)


# plays all games in a given matchday, using the provided input for how the teams
# were right at the beginning of that matchday
def computeStatesAtMatchday(matchday, prevSkills, tactics, teamOrders, matchdaySeed):
    nTeams = len(prevSkills)
    nMatchesPerMatchday = nTeams//2
    scores = np.zeros([nMatchesPerMatchday, 2], int)
    skillsAtMatchday = createEmptyPlayerStatesForAllTeams(nTeams)

    for match in range(nMatchesPerMatchday):
        team1, team2 = getTeamsInMatch(matchday, match, nTeams)

        goals1, goals2 = playMatch(
            prevSkills[team1],
            prevSkills[team2],
            tactics[team1],
            tactics[team2],
            teamOrders[team1],
            teamOrders[team2],
            matchdaySeed
        )
        scores[match] = [goals1, goals2]
        skillsAtMatchday[team1], skillsAtMatchday[team2] = \
            updatePlayerSkillsAfterMatch(
                    prevSkills[team1],
                    prevSkills[team2],
                    goals1,
                    goals2
                )
    return skillsAtMatchday, scores

# checks if 2 structs are equal by comparing the hash of their serialization
def areEqualStructs(st1, st2):
    return serialHash(st1) == serialHash(st2)


# the dataToChallenge should contain only one entry, which is the value of the first leaf.
# we just check that this leaf is the actual playerState we're testing
def isPlayerStateInsideDataToChallenge(playerState, dataToChallengePlayerState):
    return areEqualStructs(
        playerState,
        list(dataToChallengePlayerState.values.values())[0]
    )


def createEmptyPlayerStatesForAllTeams(nTeams):
    return  [[None for playerPosInLeague in range(NPLAYERS_PER_TEAM)] for team in range(nTeams)]


# ---------------- FUNCTIONS TO ADVANCE BLOCKS IN THE BC AND CLIENT ----------------
# advances both BC and CLIENT, and commits the userActions if it goes through a verse
def advanceToBlock(n, ST, ST_CLIENT):
    verseWasCrossedBC       = ST.advanceToBlock(n)
    verseWasCrossedCLIENT   = ST_CLIENT.advanceToBlock(n)
    assert verseWasCrossedBC == verseWasCrossedCLIENT, "CLIENT and BC not synced in verse crossing"
    if verseWasCrossedBC:
        ST_CLIENT.syncActions(ST)

def advanceNBlocks(deltaN, ST, ST_CLIENT):
    advanceToBlock(
        ST.currentBlock + deltaN,
        ST,
        ST_CLIENT
    )

def advanceNVerses(nVerses, ST, ST_CLIENT):
    for verse in range(nVerses):
        advanceToBlock(ST.nextVerseBlock(), ST, ST_CLIENT)


# ------------------------------------------------


# selects a pseudo-random element of an array (obtained from hashing the seed)
def getRandomElement(arr, seed):
    nElems = len(arr)
    return arr[intHash(serialize2str(seed)) % nElems]



# Merkle proof: given a tree, and its leafs,
# it creates the hashes required to prove that a given idx in the leave belongs to the tree.
# "values" is just the pair [ leafIdx, leafValue ]
def prepareProofForIdxs(idxsToProve, tree, leafs):
    neededHashes = proof(tree, idxsToProve)
    values = {}
    for leafIdx in idxsToProve:
        values[leafIdx] = leafs[leafIdx]
    return neededHashes, values

# returns an 1D-array from a 2D-array
def flatten(statesPerTeam):
    flatStates = []
    for statesTeam in statesPerTeam:
        for statePlayer in statesTeam:
            flatStates.append(MinimalPlayerState(statePlayer)) # select only skills and playerIdx
    return flatStates


# MAIN function to be called by anyone who want to make sure that the playerState is the TRULY LATEST STATE in the game
# It uses pre-hash data from CLIENT, and compares against whatever is needed in the BC
def certifyPlayerState(playerState, ST, ST_CLIENT):
    # As always we first derive the latest skills (from the last league played):
    playerChallengeData = ST_CLIENT.computeDataToChallengePlayerSkills(playerState.getPlayerIdx())
    # ...and then we update with whatever sales took place afterwards
    # TONI TODO -- update with skills!
    playerChallengeDataUpdated = ST_CLIENT.updateChallengeDataAfterLastLeaguePlayed(playerChallengeData)
    assert ST.areLatestSkills(playerState, playerChallengeDataUpdated), "Computed player state by CLIENT is not recognized by BC.."




# It uses the CLIENT data to submit a challenge to the BC
def challengeLeagueAtSelectedMatchday(selectedMatchday, leagueIdx, ST, ST_CLIENT):
    assert ST.leagues[leagueIdx].hasLeagueBeenUpdated(), "Cannot challenge a league that has not been updated"
    assert not ST.isFullyVerified(leagueIdx), "Cannot challenge a league after challenging period is over"

    # ...first, it selects a matchday, and gathers the data at that matchday (states, tactics, teamOrders)
    dataAtPrevMatchday = ST_CLIENT.getPrevMatchdayData(leagueIdx, selectedMatchday)
    # ...next, it builds the Merkle proof for the actions commited on the corresponding verse, for that league
    merkleProofDataForMatchday = ST_CLIENT.getMerkleProof(leagueIdx, selectedMatchday)

    # ...finally, it does the challenge. If successful, it will reset() the leauge update
    ST.challengeMatchdayStates(
        leagueIdx,
        selectedMatchday,
        dataAtPrevMatchday,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        duplicate(ST_CLIENT.leagues[leagueIdx].actionsPerMatchday[selectedMatchday]),
        merkleProofDataForMatchday
    )