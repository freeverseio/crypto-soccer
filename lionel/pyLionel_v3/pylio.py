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
        playerState = ST_CLIENT.getLastWrittenInClientPlayerStateFromPlayerIdx(playerIdx)
        playerChallengeData = ST_CLIENT.computeDataToChallengePlayerIdx(playerState.getPlayerIdx())
        assert ST_CLIENT.isCorrectStateForPlayerIdx(playerState, playerChallengeData), "Player state not correctly in sync"
        hash += printPlayer(playerState)
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


def playMatch(initPlayerStates1, initPlayerStates2, tactics1, tactics2, teamOrders1, teamOrders2, MatchSeed):
    hash1 = intHash(str(MatchSeed)+serialize2str(initPlayerStates1)+serialize2str(tactics1)+serialize2str(teamOrders1))
    hash2 = intHash(str(MatchSeed)+serialize2str(initPlayerStates2)+serialize2str(tactics2)+serialize2str(teamOrders2))
    return hash1 % 4, hash2 % 4

def computeTeamRating(playerStates):
    return sum([sum(thisPlayerState.getSkills()) for thisPlayerState in playerStates])

def addFixedPointsToAllPlayers(playerStates, points):
    for playerState in playerStates:
        playerState.setSkills(playerState.getSkills() + points)

def updatePlayerStatesAfterMatch(playerStates1, playerStates2, goals1, goals2):
    ps1 = duplicate(playerStates1)
    ps2 = duplicate(playerStates2)

    if goals1 == goals2:
        return ps1, ps2

    pointsWon = computePointsWon(ps1, ps2, goals1, goals2)
    if goals1 > goals2:
        addFixedPointsToAllPlayers(ps1, pointsWon)
    else:
        addFixedPointsToAllPlayers(ps2, pointsWon)

    return ps1, ps2


def computePointsWon(playerState1, playerState2, goals1, goals2):
    ratingDiff              = computeTeamRating(playerState1) - computeTeamRating(playerState2)
    winnerWasBetter         = (ratingDiff > 0 and goals1>goals2) or (ratingDiff < 0 and goals1<goals2)

    if ratingDiff == 0:
        return 5
    else:
        return (2 if winnerWasBetter else 10)


def computeStatesAtMatchday(matchday, prevStates, tactics, teamOrders, matchdaySeed):
    nTeams = len(prevStates)
    nMatchesPerMatchday = nTeams//2
    scores = np.zeros([nMatchesPerMatchday, 2], int)
    statesAtMatchday = createEmptyPlayerStatesForAllTeams(nTeams)

    for match in range(nMatchesPerMatchday):
        team1, team2 = getTeamsInMatch(matchday, match, nTeams)

        goals1, goals2 = playMatch(
            prevStates[team1],
            prevStates[team2],
            tactics[team1],
            tactics[team2],
            teamOrders[team1],
            teamOrders[team2],
            matchdaySeed
        )
        scores[match] = [goals1, goals2]
        statesAtMatchday[team1], statesAtMatchday[team2] = \
            updatePlayerStatesAfterMatch(
                    prevStates[team1],
                    prevStates[team2],
                    goals1,
                    goals2
                )
    return statesAtMatchday, scores


def getMatchsPlayerByTeam(selectedTeam, nTeams):
    matchdayMatch = []
    nMatchdays = 2*(nTeams-1)
    nMatchesPerMatchday = nTeams/2
    for matchday in range(nMatchdays):
        for match in range(nMatchesPerMatchday):
            team1, team2 = getTeamsInMatch(matchday, match, nTeams)
            if (team1==selectedTeam) or (team2==selectedTeam):
                matchdayMatch.append([matchday,match])
    return matchdayAndMatch

def areEqualStructs(st1, st2):
    return serialHash(st1) == serialHash(st2)


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


def getRandomElement(arr, seed):
    nElems = len(arr)
    return arr[intHash(serialize2str(seed)) % nElems]

# ------------------------------------------------

# Merkle proof: given a tree, and its leafs,
# it creates the hashes required to prove that a given idx in the leave belongs to the tree.
# "values" is just the pair [ leafIdx, leafValue ]
def prepareProofForIdxs(idxsToProve, tree, leafs):
    neededHashes = proof(tree, idxsToProve)
    values = {}
    for leafIdx in idxsToProve:
        values[leafIdx] = leafs[leafIdx]
    return neededHashes, values

def flatten(statesPerTeam):
    flatStates = []
    for statesTeam in statesPerTeam:
        for statePlayer in statesTeam:
            flatStates.append(statePlayer)
    return flatStates

